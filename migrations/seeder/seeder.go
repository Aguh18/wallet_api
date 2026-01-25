package seeder

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db            *gorm.DB
	userSeeder    *UserSeeder
	// Add other seeders here when needed
	// accountSeeder *AccountSeeder
	// transactionSeeder *TransactionSeeder
}

// New creates new seeder
func New(db *gorm.DB) *Seeder {
	return &Seeder{
		db:         db,
		userSeeder: NewUserSeeder(db),
	}
}

// Seed runs all seeders
func (s *Seeder) Seed(ctx context.Context) error {
	fmt.Println("🌱 Starting database seeding...")

	// Seed users
	if err := s.userSeeder.Seed(ctx); err != nil {
		return fmt.Errorf("user seeding failed: %w", err)
	}

	// Add more seeders here
	// if err := s.accountSeeder.Seed(ctx); err != nil {
	// 	return fmt.Errorf("account seeding failed: %w", err)
	// }

	fmt.Println("✅ Database seeding completed!")
	return nil
}
