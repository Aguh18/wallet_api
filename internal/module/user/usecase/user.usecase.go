package userusecase

import (
	"context"
	"fmt"

	"wallet_api/internal/common/errors"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/user/repository"
	"wallet_api/internal/utils"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, username, password string) (*entity.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateProfile(ctx context.Context, user *entity.User) error
}

type useCase struct {
	repo repository.UserRepository
}

func New(repo repository.UserRepository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (uc *useCase) Register(ctx context.Context, user *entity.User) error {
	// Check if username exists
	existing, err := uc.repo.FindByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return errors.New(409, "Username already registered", nil)
	}

	// Check if email exists
	existingEmail, err := uc.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingEmail != nil {
		return errors.New(409, "Email already registered", nil)
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

func (uc *useCase) Login(ctx context.Context, username, password string) (*entity.User, error) {
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

func (uc *useCase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.ErrNotFound
	}

	return user, nil
}

func (uc *useCase) UpdateProfile(ctx context.Context, user *entity.User) error {
	// Check if user exists
	existing, err := uc.repo.FindByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if existing == nil {
		return errors.ErrNotFound
	}

	// Check if username is taken by another user
	if user.Username != existing.Username {
		usernameTaken, err := uc.repo.FindByUsername(ctx, user.Username)
		if err != nil {
			return fmt.Errorf("failed to check username: %w", err)
		}
		if usernameTaken != nil && usernameTaken.ID != user.ID {
			return errors.New(409, "Username already taken", nil)
		}
	}

	// Check if email is taken by another user
	if user.Email != existing.Email {
		emailTaken, err := uc.repo.FindByEmail(ctx, user.Email)
		if err != nil {
			return fmt.Errorf("failed to check email: %w", err)
		}
		if emailTaken != nil && emailTaken.ID != user.ID {
			return errors.New(409, "Email already taken", nil)
		}
	}

	// Update user
	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
