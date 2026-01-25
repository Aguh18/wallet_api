package user

import (
	"wallet_api/internal/module/user/handler"
	"wallet_api/internal/module/user/repository"
	"wallet_api/internal/module/user/usecase"
	"wallet_api/pkg/logger"
	"gorm.io/gorm"
)

type Module struct {
	UseCase usecase.UseCase
	Handler *handler.Handler
}

func NewModule(db *gorm.DB, log logger.Interface) *Module {
	repo := repository.New(db)

	uc := usecase.New(repo)

	h := handler.New(uc, log)

	return &Module{
		UseCase: uc,
		Handler: h,
	}
}
