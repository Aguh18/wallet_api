package handler

import (
	"wallet_api/internal/common/response"
	"wallet_api/internal/module/account/dto/request"
	resp "wallet_api/internal/module/account/dto/response"
	accountusecase "wallet_api/internal/module/account/usecase"
	"wallet_api/pkg/logger"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	uc  accountusecase.UseCase
	log logger.Interface
}

func New(uc accountusecase.UseCase, log logger.Interface) *Handler {
	return &Handler{
		uc:  uc,
		log: log,
	}
}

func (h *Handler) CreateAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	req := new(request.CreateAccountRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	account, err := h.uc.CreateAccount(c.Context(), userID, req.AccountName, req.Currency)
	if err != nil {
		h.log.Error("failed to create account: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to create account"))
	}

	return c.JSON(response.Success(resp.ToAccountDto(account), "Account created successfully"))
}

func (h *Handler) GetAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accountID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid account ID"))
	}

	account, err := h.uc.GetAccount(c.Context(), accountID)
	if err != nil {
		h.log.Error("failed to get account: %v", err)
		return c.Status(404).JSON(response.Error(404, "Account not found"))
	}

	return c.JSON(response.Success(resp.ToAccountDto(account), "Account retrieved"))
}

func (h *Handler) GetUserAccounts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	accounts, err := h.uc.GetUserAccounts(c.Context(), userID)
	if err != nil {
		h.log.Error("failed to get user accounts: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to get accounts"))
	}

	return c.JSON(response.Success(resp.ToAccountDtos(accounts), "Accounts retrieved"))
}

func (h *Handler) Deposit(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accountID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid account ID"))
	}

	req := new(request.TransactionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	if err := h.uc.Deposit(c.Context(), accountID, req.Amount, req.Description); err != nil {
		h.log.Error("failed to deposit: %v", err)
		return c.Status(400).JSON(response.Error(400, err.Error()))
	}

	return c.JSON(response.Success(nil, "Deposit successful"))
}

func (h *Handler) Withdraw(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accountID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid account ID"))
	}

	req := new(request.TransactionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	if err := h.uc.Withdraw(c.Context(), accountID, req.Amount, req.Description); err != nil {
		h.log.Error("failed to withdraw: %v", err)
		return c.Status(400).JSON(response.Error(400, err.Error()))
	}

	return c.JSON(response.Success(nil, "Withdrawal successful"))
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	idParam := c.Params("id")
	accountID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid account ID"))
	}

	limit := 10
	offset := 0

	if l := c.QueryInt("limit", 10); l > 0 {
		limit = l
	}
	if o := c.QueryInt("offset", 0); o >= 0 {
		offset = o
	}

	transactions, err := h.uc.GetTransactions(c.Context(), accountID, limit, offset)
	if err != nil {
		h.log.Error("failed to get transactions: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to get transactions"))
	}

	return c.JSON(response.Success(resp.ToTransactionDtos(transactions), "Transactions retrieved"))
}
