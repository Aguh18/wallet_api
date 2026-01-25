package handler

import (
	"wallet_api/internal/common/response"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/user/usecase"
	"wallet_api/pkg/logger"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	uc  *usecase.UseCase
	log logger.Interface
}

// New creates new user handler
func New(uc *usecase.UseCase, log logger.Interface) *Handler {
	return &Handler{
		uc:  uc,
		log: log,
	}
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

// Register handles user registration
func (h *Handler) Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	user := &entity.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	if err := h.uc.Register(c.Context(), user); err != nil {
		h.log.Error("failed to register user: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to register user"))
	}

	user.PasswordHash = ""

	return c.JSON(response.Success(user, "User registered successfully"))
}

// Login handles user login
func (h *Handler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	user, err := h.uc.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		h.log.Error("failed to login: %v", err)
		return c.Status(401).JSON(response.Error(401, "Invalid credentials"))
	}

	user.PasswordHash = ""

	return c.JSON(response.Success(user, "Login successful"))
}

// GetProfile handles get user profile
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	user, err := h.uc.GetProfile(c.Context(), userID)
	if err != nil {
		h.log.Error("failed to get profile: %v", err)
		return c.Status(404).JSON(response.Error(404, "User not found"))
	}

	return c.JSON(response.Success(user, "Profile retrieved"))
}

// UpdateProfile handles update user profile
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	req := new(UpdateProfileRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	user, err := h.uc.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(404).JSON(response.Error(404, "User not found"))
	}

	user.Username = req.Username

	if err := h.uc.UpdateProfile(c.Context(), user); err != nil {
		h.log.Error("failed to update profile: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to update profile"))
	}

	return c.JSON(response.Success(user, "Profile updated"))
}
