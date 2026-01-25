// Package postgres implements postgres connection using GORM.
package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *gorm.DB
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	var err error
	for pg.connAttempts > 0 {
		pg.DB, err = gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	sqlDB, err := pg.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres - GetDB: %w", err)
	}

	sqlDB.SetMaxOpenConns(pg.maxPoolSize)
	sqlDB.SetMaxIdleConns(pg.maxPoolSize)

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.DB != nil {
		sqlDB, err := p.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
