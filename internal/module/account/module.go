package account

import (
	"wallet_api/internal/module/account/handler"
	accountrepository "wallet_api/internal/module/account/repository"
	"wallet_api/internal/module/account/usecase"
	"wallet_api/pkg/logger"
	"gorm.io/gorm"
)

// Module represents account module with all its components
type Module struct {
	Repository accountrepository.Repository
	UseCase    *usecase.UseCase
	Handler    *handler.Handler
}

// NewModule creates and initializes account module
func NewModule(db *gorm.DB, log logger.Interface) *Module {
	repo := accountrepository.New(db)
	uc := usecase.New(repo)
	h := handler.New(uc, log)

	return &Module{
		Repository: repo,
		UseCase:    uc,
		Handler:    h,
	}
}
