package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	// Base settings
	host     = "localhost"
	attempts = 20

	// Connection settings
	httpURL        = "http://" + host + ":8000"
	healthPath     = httpURL + "/healthz"
	requestTimeout = 10 * time.Second

	// API paths
	basePathV1        = httpURL + "/v1"
	authPath          = basePathV1 + "/auth"
	userPath          = basePathV1 + "/users"
	accountPath       = basePathV1 + "/accounts"
)

var (
	testHTTPClient *http.Client
)

// TestResponse represents standard API response
type TestResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *TestError  `json:"error,omitempty"`
}

type TestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// User data structures
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// Account data structures
type CreateAccountRequest struct {
	AccountName string `json:"account_name"`
	Currency    string `json:"currency"`
}

type AccountResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	AccountName string `json:"account_name"`
	Currency    string `json:"currency"`
	Balance     int64  `json:"balance"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type TransactionRequest struct {
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type TransferRequest struct {
	ToAccountID string `json:"to_account_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type TransactionResponse struct {
	ID            string `json:"id"`
	AccountID     string `json:"account_id"`
	ReferenceID   string `json:"reference_id"`
	Type          string `json:"type"`
	Amount        int64  `json:"amount"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
}

func init() {
	testHTTPClient = &http.Client{
		Timeout: requestTimeout,
	}
}

// Helper function to make HTTP requests
func makeRequest(method, url string, body interface{}, cookies []*http.Cookie) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add cookies if provided
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return testHTTPClient.Do(req)
}

// Helper function to parse response
func parseResponse(resp *http.Response) (*TestResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	var testResp TestResponse
	if err := json.Unmarshal(body, &testResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &testResp, nil
}

// Health check
func getHealthCheck(url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return -1, err
	}

	resp, err := testHTTPClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func healthCheck(attempts int) error {
	for attempts > 0 {
		statusCode, err := getHealthCheck(healthPath)
		if err != nil {
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)
		time.Sleep(time.Second)
		attempts--
	}

	return fmt.Errorf("url %s is not available", healthPath)
}

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: %s", err)
	}

	log.Printf("Integration tests: httpURL %s is available", httpURL)

	code := m.Run()
	os.Exit(code)
}

// ============================================================================
// AUTHENTICATION TESTS
// ============================================================================

func TestUserRegistrationAndLogin(t *testing.T) {
	// Generate unique username
	testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])

	t.Run("Register New User", func(t *testing.T) {
		req := RegisterRequest{
			Username: testUsername,
			Password: "password123",
		}

		resp, err := makeRequest(http.MethodPost, authPath+"/register", req, nil)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Check if auth cookies are set
		cookies := resp.Cookies()
		hasAccessToken := false
		hasRefreshToken := false

		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				hasAccessToken = true
			}
			if cookie.Name == "refresh_token" {
				hasRefreshToken = true
			}
		}

		if !hasAccessToken || !hasRefreshToken {
			t.Error("Auth cookies not set after registration")
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}
	})

	t.Run("Register Duplicate User Should Fail", func(t *testing.T) {
		req := RegisterRequest{
			Username: testUsername,
			Password: "password123",
		}

		resp, err := makeRequest(http.MethodPost, authPath+"/register", req, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("Login with Valid Credentials", func(t *testing.T) {
		req := LoginRequest{
			Username: testUsername,
			Password: "password123",
		}

		resp, err := makeRequest(http.MethodPost, authPath+"/login", req, nil)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}
	})

	t.Run("Login with Invalid Credentials", func(t *testing.T) {
		req := LoginRequest{
			Username: testUsername,
			Password: "wrongpassword",
		}

		resp, err := makeRequest(http.MethodPost, authPath+"/login", req, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestUserProfile(t *testing.T) {
	// Register and login a new user
	testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])

	registerReq := RegisterRequest{
		Username: testUsername,
		Password: "password123",
	}

	resp, err := makeRequest(http.MethodPost, authPath+"/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to register user, status: %d", resp.StatusCode)
	}

	cookies := resp.Cookies()

	t.Run("Get User Profile", func(t *testing.T) {
		resp, err := makeRequest(http.MethodGet, userPath+"/profile", nil, cookies)
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Parse user data
		data, err := json.Marshal(testResp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var user UserResponse
		if err := json.Unmarshal(data, &user); err != nil {
			t.Fatalf("Failed to unmarshal user: %v", err)
		}

		if user.Username != testUsername {
			t.Errorf("Expected username %s, got %s", testUsername, user.Username)
		}
	})

	t.Run("Update User Profile", func(t *testing.T) {
		newUsername := fmt.Sprintf("updated_%s", uuid.New().String()[:8])
		updateReq := UpdateProfileRequest{
			Username: newUsername,
		}

		resp, err := makeRequest(http.MethodPut, userPath+"/profile", updateReq, cookies)
		if err != nil {
			t.Fatalf("Failed to update profile: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Verify the update
		data, err := json.Marshal(testResp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var user UserResponse
		if err := json.Unmarshal(data, &user); err != nil {
			t.Fatalf("Failed to unmarshal user: %v", err)
		}

		if user.Username != newUsername {
			t.Errorf("Expected username %s, got %s", newUsername, user.Username)
		}
	})

	t.Run("Access Profile Without Auth Should Fail", func(t *testing.T) {
		resp, err := makeRequest(http.MethodGet, userPath+"/profile", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestLogoutAndRefresh(t *testing.T) {
	testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])

	// Register user
	registerReq := RegisterRequest{
		Username: testUsername,
		Password: "password123",
	}

	resp, err := makeRequest(http.MethodPost, authPath+"/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	cookies := resp.Cookies()

	t.Run("Logout User", func(t *testing.T) {
		resp, err := makeRequest(http.MethodPost, authPath+"/logout", nil, cookies)
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}
	})

	t.Run("Refresh Token", func(t *testing.T) {
		// Login again to get fresh tokens
		loginReq := LoginRequest{
			Username: testUsername,
			Password: "password123",
		}

		resp, err := makeRequest(http.MethodPost, authPath+"/login", loginReq, nil)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		cookies := resp.Cookies()

		// Use refresh token to get new access token
		resp, err = makeRequest(http.MethodPost, authPath+"/refresh", nil, cookies)
		if err != nil {
			t.Fatalf("Failed to refresh token: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}
	})
}

// ============================================================================
// ACCOUNT TESTS
// ============================================================================

func TestAccountManagement(t *testing.T) {
	// Register and login a new user
	testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])

	registerReq := RegisterRequest{
		Username: testUsername,
		Password: "password123",
	}

	resp, err := makeRequest(http.MethodPost, authPath+"/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to register user, status: %d", resp.StatusCode)
	}

	cookies := resp.Cookies()

	var accountID string

	t.Run("Create Account", func(t *testing.T) {
		createReq := CreateAccountRequest{
			AccountName: "My Wallet",
			Currency:    "IDR",
		}

		resp, err := makeRequest(http.MethodPost, accountPath, createReq, cookies)
		if err != nil {
			t.Fatalf("Failed to create account: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Parse account data
		data, err := json.Marshal(testResp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var account AccountResponse
		if err := json.Unmarshal(data, &account); err != nil {
			t.Fatalf("Failed to unmarshal account: %v", err)
		}

		accountID = account.ID

		if account.AccountName != "My Wallet" {
			t.Errorf("Expected account name 'My Wallet', got %s", account.AccountName)
		}

		if account.Balance != 0 {
			t.Errorf("Expected initial balance 0, got %d", account.Balance)
		}

		if account.Status != "active" {
			t.Errorf("Expected status 'active', got %s", account.Status)
		}
	})

	t.Run("Get Account by ID", func(t *testing.T) {
		if accountID == "" {
			t.Skip("Account ID not available")
		}

		url := fmt.Sprintf("%s/%s", accountPath, accountID)
		resp, err := makeRequest(http.MethodGet, url, nil, cookies)
		if err != nil {
			t.Fatalf("Failed to get account: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}
	})

	t.Run("Get User Accounts", func(t *testing.T) {
		resp, err := makeRequest(http.MethodGet, accountPath, nil, cookies)
		if err != nil {
			t.Fatalf("Failed to get user accounts: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Parse accounts data
		data, err := json.Marshal(testResp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var accounts []AccountResponse
		if err := json.Unmarshal(data, &accounts); err != nil {
			t.Fatalf("Failed to unmarshal accounts: %v", err)
		}

		if len(accounts) == 0 {
			t.Error("Expected at least one account")
		}
	})

	t.Run("Get Non-Existent Account Should Fail", func(t *testing.T) {
		fakeID := uuid.New().String()
		url := fmt.Sprintf("%s/%s", accountPath, fakeID)
		resp, err := makeRequest(http.MethodGet, url, nil, cookies)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

// ============================================================================
// TRANSACTION TESTS
// ============================================================================

func TestTransactions(t *testing.T) {
	// Register and login a new user
	testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])

	registerReq := RegisterRequest{
		Username: testUsername,
		Password: "password123",
	}

	resp, err := makeRequest(http.MethodPost, authPath+"/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	cookies := resp.Cookies()

	// Create account
	createReq := CreateAccountRequest{
		AccountName: "Test Wallet",
		Currency:    "IDR",
	}

	resp, err = makeRequest(http.MethodPost, accountPath, createReq, cookies)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	testResp, err := parseResponse(resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, err := json.Marshal(testResp.Data)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var account AccountResponse
	if err := json.Unmarshal(data, &account); err != nil {
		t.Fatalf("Failed to unmarshal account: %v", err)
	}

	accountID := account.ID

	t.Run("Deposit Funds", func(t *testing.T) {
		depositReq := TransactionRequest{
			Amount:      100000,
			Description: "Initial deposit",
		}

		url := fmt.Sprintf("%s/%s/deposit", accountPath, accountID)
		resp, err := makeRequest(http.MethodPost, url, depositReq, cookies)
		if err != nil {
			t.Fatalf("Failed to deposit: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Verify balance updated
		url = fmt.Sprintf("%s/%s", accountPath, accountID)
		resp, err = makeRequest(http.MethodGet, url, nil, cookies)
		if err != nil {
			t.Fatalf("Failed to get account: %v", err)
		}

		testResp, err = parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		data, _ = json.Marshal(testResp.Data)
		var updatedAccount AccountResponse
		if err := json.Unmarshal(data, &updatedAccount); err != nil {
			t.Fatalf("Failed to unmarshal account: %v", err)
		}

		if updatedAccount.Balance != 100000 {
			t.Errorf("Expected balance 100000, got %d", updatedAccount.Balance)
		}
	})

	t.Run("Withdraw Funds", func(t *testing.T) {
		withdrawReq := TransactionRequest{
			Amount:      50000,
			Description: "Cash withdrawal",
		}

		url := fmt.Sprintf("%s/%s/withdraw", accountPath, accountID)
		resp, err := makeRequest(http.MethodPost, url, withdrawReq, cookies)
		if err != nil {
			t.Fatalf("Failed to withdraw: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Verify balance updated
		url = fmt.Sprintf("%s/%s", accountPath, accountID)
		resp, err = makeRequest(http.MethodGet, url, nil, cookies)
		if err != nil {
			t.Fatalf("Failed to get account: %v", err)
		}

		testResp, err = parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		data, _ = json.Marshal(testResp.Data)
		var updatedAccount AccountResponse
		if err := json.Unmarshal(data, &updatedAccount); err != nil {
			t.Fatalf("Failed to unmarshal account: %v", err)
		}

		if updatedAccount.Balance != 50000 {
			t.Errorf("Expected balance 50000, got %d", updatedAccount.Balance)
		}
	})

	t.Run("Withdraw More Than Balance Should Fail", func(t *testing.T) {
		withdrawReq := TransactionRequest{
			Amount:      100000,
			Description: "Overdraft attempt",
		}

		url := fmt.Sprintf("%s/%s/withdraw", accountPath, accountID)
		resp, err := makeRequest(http.MethodPost, url, withdrawReq, cookies)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Get Transaction History", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s/transactions", accountPath, accountID)
		resp, err := makeRequest(http.MethodGet, url, nil, cookies)
		if err != nil {
			t.Fatalf("Failed to get transactions: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Parse transactions data
		data, err := json.Marshal(testResp.Data)
		if err != nil {
			t.Fatalf("Failed to marshal data: %v", err)
		}

		var transactions []TransactionResponse
		if err := json.Unmarshal(data, &transactions); err != nil {
			t.Fatalf("Failed to unmarshal transactions: %v", err)
		}

		// Should have at least 2 transactions (deposit and withdraw)
		if len(transactions) < 2 {
			t.Errorf("Expected at least 2 transactions, got %d", len(transactions))
		}
	})
}

func TestTransferBetweenAccounts(t *testing.T) {
	// Create first user
	user1Username := fmt.Sprintf("testuser1_%s", uuid.New().String()[:8])

	registerReq1 := RegisterRequest{
		Username: user1Username,
		Password: "password123",
	}

	resp, err := makeRequest(http.MethodPost, authPath+"/register", registerReq1, nil)
	if err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}

	cookies1 := resp.Cookies()

	// Create account for user1
	createReq1 := CreateAccountRequest{
		AccountName: "User1 Wallet",
		Currency:    "IDR",
	}

	resp, err = makeRequest(http.MethodPost, accountPath, createReq1, cookies1)
	if err != nil {
		t.Fatalf("Failed to create account for user1: %v", err)
	}

	testResp, _ := parseResponse(resp)
	data, _ := json.Marshal(testResp.Data)
	var account1 AccountResponse
	if err := json.Unmarshal(data, &account1); err != nil {
		t.Fatalf("Failed to unmarshal account1: %v", err)
	}

	// Deposit to user1 account
	depositReq := TransactionRequest{
		Amount:      200000,
		Description: "Initial deposit",
	}

	url := fmt.Sprintf("%s/%s/deposit", accountPath, account1.ID)
	resp, err = makeRequest(http.MethodPost, url, depositReq, cookies1)
	if err != nil {
		t.Fatalf("Failed to deposit to user1: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to deposit, status: %d", resp.StatusCode)
	}

	// Create second user
	user2Username := fmt.Sprintf("testuser2_%s", uuid.New().String()[:8])

	registerReq2 := RegisterRequest{
		Username: user2Username,
		Password: "password123",
	}

	resp, err = makeRequest(http.MethodPost, authPath+"/register", registerReq2, nil)
	if err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}

	cookies2 := resp.Cookies()

	// Create account for user2
	createReq2 := CreateAccountRequest{
		AccountName: "User2 Wallet",
		Currency:    "IDR",
	}

	resp, err = makeRequest(http.MethodPost, accountPath, createReq2, cookies2)
	if err != nil {
		t.Fatalf("Failed to create account for user2: %v", err)
	}

	testResp, _ = parseResponse(resp)
	data, _ = json.Marshal(testResp.Data)
	var account2 AccountResponse
	if err := json.Unmarshal(data, &account2); err != nil {
		t.Fatalf("Failed to unmarshal account2: %v", err)
	}

	t.Run("Transfer Between Accounts", func(t *testing.T) {
		transferReq := TransferRequest{
			ToAccountID: account2.ID,
			Amount:      100000,
			Description: "Transfer to user2",
		}

		url := fmt.Sprintf("%s/%s/transfer", accountPath, account1.ID)
		resp, err := makeRequest(http.MethodPost, url, transferReq, cookies1)
		if err != nil {
			t.Fatalf("Failed to transfer: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		testResp, err := parseResponse(resp)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !testResp.Success {
			t.Errorf("Expected success=true, got false: %s", testResp.Message)
		}

		// Verify user1 balance decreased
		url = fmt.Sprintf("%s/%s", accountPath, account1.ID)
		resp, err = makeRequest(http.MethodGet, url, nil, cookies1)
		if err != nil {
			t.Fatalf("Failed to get user1 account: %v", err)
		}

		testResp, _ = parseResponse(resp)
		data, _ = json.Marshal(testResp.Data)
		var updatedAccount1 AccountResponse
		if err := json.Unmarshal(data, &updatedAccount1); err != nil {
			t.Fatalf("Failed to unmarshal updated account1: %v", err)
		}

		if updatedAccount1.Balance != 100000 {
			t.Errorf("Expected user1 balance 100000, got %d", updatedAccount1.Balance)
		}

		// Verify user2 balance increased
		url = fmt.Sprintf("%s/%s", accountPath, account2.ID)
		resp, err = makeRequest(http.MethodGet, url, nil, cookies2)
		if err != nil {
			t.Fatalf("Failed to get user2 account: %v", err)
		}

		testResp, _ = parseResponse(resp)
		data, _ = json.Marshal(testResp.Data)
		var updatedAccount2 AccountResponse
		if err := json.Unmarshal(data, &updatedAccount2); err != nil {
			t.Fatalf("Failed to unmarshal updated account2: %v", err)
		}

		if updatedAccount2.Balance != 100000 {
			t.Errorf("Expected user2 balance 100000, got %d", updatedAccount2.Balance)
		}
	})

	t.Run("Transfer More Than Balance Should Fail", func(t *testing.T) {
		transferReq := TransferRequest{
			ToAccountID: account2.ID,
			Amount:      500000,
			Description: "Overdraft transfer",
		}

		url := fmt.Sprintf("%s/%s/transfer", accountPath, account1.ID)
		resp, err := makeRequest(http.MethodPost, url, transferReq, cookies1)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Transfer to Non-Existent Account Should Fail", func(t *testing.T) {
		fakeID := uuid.New().String()
		transferReq := TransferRequest{
			ToAccountID: fakeID,
			Amount:      10000,
			Description: "Transfer to fake account",
		}

		url := fmt.Sprintf("%s/%s/transfer", accountPath, account1.ID)
		resp, err := makeRequest(http.MethodPost, url, transferReq, cookies1)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

// ============================================================================
// EDGE CASES AND VALIDATION TESTS
// ============================================================================

func TestEdgeCases(t *testing.T) {
	testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])

	registerReq := RegisterRequest{
		Username: testUsername,
		Password: "password123",
	}

	resp, err := makeRequest(http.MethodPost, authPath+"/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	cookies := resp.Cookies()

	// Create account
	createReq := CreateAccountRequest{
		AccountName: "Test Wallet",
		Currency:    "IDR",
	}

	resp, err = makeRequest(http.MethodPost, accountPath, createReq, cookies)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	testResp, _ := parseResponse(resp)
	data, _ := json.Marshal(testResp.Data)
	var account AccountResponse
	if err := json.Unmarshal(data, &account); err != nil {
		t.Fatalf("Failed to unmarshal account: %v", err)
	}

	t.Run("Deposit with Negative Amount Should Fail", func(t *testing.T) {
		depositReq := TransactionRequest{
			Amount:      -1000,
			Description: "Negative deposit",
		}

		url := fmt.Sprintf("%s/%s/deposit", accountPath, account.ID)
		resp, err := makeRequest(http.MethodPost, url, depositReq, cookies)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Withdraw with Negative Amount Should Fail", func(t *testing.T) {
		withdrawReq := TransactionRequest{
			Amount:      -1000,
			Description: "Negative withdrawal",
		}

		url := fmt.Sprintf("%s/%s/withdraw", accountPath, account.ID)
		resp, err := makeRequest(http.MethodPost, url, withdrawReq, cookies)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Create Account with Invalid Currency", func(t *testing.T) {
		createReq := CreateAccountRequest{
			AccountName: "Invalid Account",
			Currency:    "XXX",
		}

		resp, err := makeRequest(http.MethodPost, accountPath, createReq, cookies)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		// This should fail validation
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected non-200 status for invalid currency")
		}
	})
}
