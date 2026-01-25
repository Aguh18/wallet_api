package account

import (
	"wallet_api/internal/module/account/handler"
	"wallet_api/internal/module/account/repository"
	accountusecase "wallet_api/internal/module/account/usecase"
	"wallet_api/pkg/logger"

	"gorm.io/gorm"
)

type Module struct {
	UseCase accountusecase.UseCase
	Handler *handler.Handler
}

func NewModule(db *gorm.DB, log logger.Interface) *Module {
	accountRepo := repository.New(db)
	transactionRepo := repository.NewTransactionRepository(db)
	uc := accountusecase.New(accountRepo, transactionRepo)
	h := handler.New(uc, log)

	return &Module{
		UseCase: uc,
		Handler: h,
	}
}
