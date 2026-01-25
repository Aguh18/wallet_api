package router

import (
	"wallet_api/internal/module/account"
	"wallet_api/internal/module/user"
	"wallet_api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Module represents all router modules
type Module struct {
	User    *user.Module
	Account *account.Module
}

// NewModule creates and initializes all router modules
func NewModule(db *gorm.DB, log logger.Interface) *Module {
	// Initialize User Module
	userModule := user.NewModule(db, log)

	// Initialize Account Module
	accountModule := account.NewModule(db, log)

	return &Module{
		User:    userModule,
		Account: accountModule,
	}
}

// RegisterRoutes registers all module routes
func (m *Module) RegisterRoutes(app *fiber.App) {
	m.User.RegisterRoutes(app)
	m.Account.RegisterRoutes(app)
}
