package user

import (
	"wallet_api/internal/module/user/handler"
	userrepository "wallet_api/internal/module/user/repository"
	"wallet_api/internal/module/user/usecase"
	"wallet_api/pkg/logger"
	"gorm.io/gorm"
)

// Module represents user module with all its components
type Module struct {
	Repository userrepository.Repository
	UseCase    *usecase.UseCase
	Handler    *handler.Handler
}

// NewModule creates and initializes user module
func NewModule(db *gorm.DB, log logger.Interface) *Module {
	// Initialize Repository
	repo := userrepository.New(db)

	// Initialize UseCase
	uc := usecase.New(repo)

	// Initialize Handler
	h := handler.New(uc, log)

	return &Module{
		Repository: repo,
		UseCase:    uc,
		Handler:    h,
	}
}
