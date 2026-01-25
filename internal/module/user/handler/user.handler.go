package handler

import (
	"time"
	"wallet_api/internal/common/errors"
	"wallet_api/internal/common/response"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/user/dto/request"
	resp "wallet_api/internal/module/user/dto/response"
	userusecase "wallet_api/internal/module/user/usecase"
	"wallet_api/internal/utils"
	"wallet_api/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	uc         userusecase.UseCase
	log        logger.Interface
	jwtManager *utils.JWTManager
}

func New(uc userusecase.UseCase, log logger.Interface) *Handler {
	return &Handler{
		uc:         uc,
		log:        log,
		jwtManager: utils.NewJWTManager(utils.GetSecretKey()),
	}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	req := new(request.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	user := &entity.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	if err := h.uc.Register(c.Context(), user); err != nil {
		h.log.Error("failed to register user: %v", err)

		// Check if it's a conflict error (user already exists)
		if appErr, ok := err.(*errors.AppError); ok {
			return c.Status(appErr.Code).JSON(response.Error(appErr.Code, appErr.Message))
		}

		return c.Status(500).JSON(response.Error(500, "Failed to register user"))
	}

	// Generate JWT tokens
	tokenPair, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		h.log.Error("failed to generate tokens: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to generate tokens"))
	}

	// Set auth cookies
	utils.SetAuthCookies(c, tokenPair.AccessToken, tokenPair.RefreshToken, time.Duration(tokenPair.ExpiresIn)*time.Second)

	return c.JSON(response.Success(resp.ToUserDto(user), "User registered successfully"))
}

func (h *Handler) Login(c *fiber.Ctx) error {
	req := new(request.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(response.Error(400, "Invalid request body"))
	}

	user, err := h.uc.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		h.log.Error("failed to login: %v", err)
		return c.Status(401).JSON(response.Error(401, "Invalid credentials"))
	}

	// Generate JWT tokens
	tokenPair, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		h.log.Error("failed to generate tokens: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to generate tokens"))
	}

	// Set auth cookies
	utils.SetAuthCookies(c, tokenPair.AccessToken, tokenPair.RefreshToken, time.Duration(tokenPair.ExpiresIn)*time.Second)

	return c.JSON(response.Success(resp.ToUserDto(user), "Login successful"))
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	user, err := h.uc.GetProfile(c.Context(), userID)
	if err != nil {
		h.log.Error("failed to get profile: %v", err)
		return c.Status(404).JSON(response.Error(404, "User not found"))
	}

	return c.JSON(response.Success(resp.ToUserDto(user), "Profile retrieved"))
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	req := new(request.UpdateProfileRequest)
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

	return c.JSON(response.Success(resp.ToUserDto(user), "Profile updated"))
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	// Clear auth cookies
	utils.ClearAuthCookies(c)

	return c.JSON(response.Success(nil, "Logout successful"))
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	// Get refresh token from cookie
	refreshToken := utils.GetRefreshTokenFromCookie(c)
	if refreshToken == "" {
		return c.Status(401).JSON(response.Error(401, "Refresh token not found"))
	}

	// Validate refresh token
	claims, err := h.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		h.log.Error("failed to validate refresh token: %v", err)
		return c.Status(401).JSON(response.Error(401, "Invalid refresh token"))
	}

	// Generate new token pair
	tokenPair, err := h.jwtManager.GenerateToken(claims.UserID, claims.Username)
	if err != nil {
		h.log.Error("failed to generate new tokens: %v", err)
		return c.Status(500).JSON(response.Error(500, "Failed to generate tokens"))
	}

	// Set new auth cookies
	utils.SetAuthCookies(c, tokenPair.AccessToken, tokenPair.RefreshToken, time.Duration(tokenPair.ExpiresIn)*time.Second)

	return c.JSON(response.Success(nil, "Token refreshed successfully"))
}
