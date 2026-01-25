package seeder

import (
	"context"
	"fmt"

	"wallet_api/internal/entity"
	"wallet_api/internal/utils"
	"gorm.io/gorm"
)

// UserSeeder handles user data seeding
type UserSeeder struct {
	db *gorm.DB
}

// NewUserSeeder creates new user seeder
func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// Seed runs the user seeding
func (s *UserSeeder) Seed(ctx context.Context) error {
	// Check if users already exist
	var count int64
	if err := s.db.WithContext(ctx).Table("users").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check users: %w", err)
	}

	if count > 0 {
		fmt.Printf("✅ Users already exist (%d records). Skipping seeding.\n", count)
		return nil
	}

	// Create sample users
	users := []entity.User{
		{
			Username:     "admin",
			PasswordHash: mustHash("admin123"),
		},
		{
			Username:     "johndoe",
			PasswordHash: mustHash("password123"),
		},
		{
			Username:     "janedoe",
			PasswordHash: mustHash("password123"),
		},
	}

	for _, user := range users {
		if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Username, err)
		}
		fmt.Printf("✅ Created user: %s\n", user.Username)
	}

	fmt.Printf("✅ User seeding completed. Total users: %d\n", len(users))
	return nil
}

// mustHash hashes password or panics
func mustHash(password string) string {
	hash, err := utils.HashPassword(password)
	if err != nil {
		panic(fmt.Sprintf("failed to hash password: %v", err))
	}
	return hash
}
