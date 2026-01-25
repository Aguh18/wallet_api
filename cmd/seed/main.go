package main

import (
	"context"
	"log"

	"wallet_api/config"
	"wallet_api/migrations/seeder"
	"wallet_api/pkg/postgres"
)

func main() {
	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Connect to database
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatalf("Database connection error: %s", err)
	}
	defer pg.Close()

	// Run seeder
	s := seeder.New(pg.DB)
	if err := s.Seed(context.Background()); err != nil {
		log.Fatalf("Seeding failed: %s", err)
	}
}
