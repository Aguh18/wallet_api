package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

type (
	// Config -.
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		PG   PG
		JWT  JWT
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// JWT -.
	JWT struct {
		Secret             string `env:"JWT_SECRET,required"`
		AccessTokenExpiry  int    `env:"ACCESS_TOKEN_EXPIRY" envDefault:"15"`
		RefreshTokenExpiry int    `env:"REFRESH_TOKEN_EXPIRY" envDefault:"7"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
