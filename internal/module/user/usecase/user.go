package usecase

import (
	"context"
	"fmt"

	"wallet_api/internal/common/errors"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/user/repository"
	"wallet_api/internal/utils"
	"github.com/google/uuid"
)

type UseCase struct {
	repo repository.Repository
}

// New creates new user usecase
func New(repo repository.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// Register creates new user
func (uc *UseCase) Register(ctx context.Context, user *entity.User) error {
	// Check if username exists
	existing, err := uc.repo.FindByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return errors.New(409, "Username already registered", nil)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashedPassword

	// Create user
	if err := uc.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Login authenticates user
func (uc *UseCase) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.ErrBadRequest
	}

	// Verify password
	if err := utils.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, errors.ErrUnauthorized
	}

	return user, nil
}

// GetProfile retrieves user by ID
func (uc *UseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.ErrNotFound
	}

	return user, nil
}

// UpdateProfile updates user profile
func (uc *UseCase) UpdateProfile(ctx context.Context, user *entity.User) error {
	// Check if user exists
	existing, err := uc.repo.FindByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if existing == nil {
		return errors.ErrNotFound
	}

	// Update user
	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
