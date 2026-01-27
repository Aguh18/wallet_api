package handler

import (
	"wallet_api/internal/common/response"
	"wallet_api/internal/module/account/dto/request"
	resp "wallet_api/internal/module/account/dto/response"
	accountusecase "wallet_api/internal/module/account/usecase"
	"wallet_api/pkg/logger"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
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

	wallet, err := h.uc.CreateWallet(c.Context(), userID, req.AccountName, req.Currency)
	if err != nil {
		h.log.Error("failed to create wallet: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to create wallet"))
	}

	return c.JSON(response.Success(resp.ToWalletDto(wallet), "Wallet created successfully"))
}

func (h *Handler) GetAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	walletID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid wallet ID"))
	}

	wallet, err := h.uc.GetWallet(c.Context(), walletID)
	if err != nil {
		h.log.Error("failed to get wallet: %v", err)
		return c.Status(404).JSON(response.Error(404, "Wallet not found"))
	}

	return c.JSON(response.Success(resp.ToWalletDto(wallet), "Wallet retrieved"))
}

func (h *Handler) GetUserAccounts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	wallets, err := h.uc.GetUserWallets(c.Context(), userID)
	if err != nil {
		h.log.Error("failed to get user wallets: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to get wallets"))
	}

	return c.JSON(response.Success(resp.ToWalletDtos(wallets), "Wallets retrieved"))
}

func (h *Handler) Deposit(c *fiber.Ctx) error {
	idParam := c.Params("id")
	walletID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid wallet ID"))
	}

	req := new(request.TransactionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	// Parse amount string to decimal
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid amount format"))
	}

	if err := h.uc.Deposit(c.Context(), walletID, amount, req.Description); err != nil {
		h.log.Error("failed to deposit: %v", err)
		return c.Status(400).JSON(response.Error(400, err.Error()))
	}

	return c.JSON(response.Success(nil, "Deposit successful"))
}

func (h *Handler) Withdraw(c *fiber.Ctx) error {
	idParam := c.Params("id")
	walletID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid wallet ID"))
	}

	req := new(request.TransactionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	// Parse amount string to decimal
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid amount format"))
	}

	if err := h.uc.Withdraw(c.Context(), walletID, amount, req.Description); err != nil {
		h.log.Error("failed to withdraw: %v", err)
		return c.Status(400).JSON(response.Error(400, err.Error()))
	}

	return c.JSON(response.Success(nil, "Withdrawal successful"))
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	idParam := c.Params("id")
	walletID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid wallet ID"))
	}

	limit := 10
	offset := 0

	if l := c.QueryInt("limit", 10); l > 0 {
		limit = l
	}
	if o := c.QueryInt("offset", 0); o >= 0 {
		offset = o
	}

	transactions, err := h.uc.GetTransactions(c.Context(), walletID, limit, offset)
	if err != nil {
		h.log.Error("failed to get transactions: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to get transactions"))
	}

	return c.JSON(response.Success(resp.ToTransactionDtos(transactions), "Transactions retrieved"))
}

func (h *Handler) Transfer(c *fiber.Ctx) error {
	fromWalletIDParam := c.Params("id")
	fromWalletID, err := uuid.Parse(fromWalletIDParam)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid from wallet ID"))
	}

	req := new(request.TransferRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	toWalletID, err := uuid.Parse(req.ToWalletID)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid to wallet ID"))
	}

	// Parse amount string to decimal
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid amount format"))
	}

	if err := h.uc.Transfer(c.Context(), fromWalletID, toWalletID, amount, req.Description); err != nil {
		h.log.Error("failed to transfer: %v", err)
		return c.Status(400).JSON(response.Error(400, err.Error()))
	}

	return c.JSON(response.Success(nil, "Transfer successful"))
}
