package account

import (
	"wallet_api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(app *fiber.App) {
	wallets := app.Group("/v1/wallets", middleware.JWTAuth())
	{

		wallets.Post("/", m.Handler.CreateAccount)
		wallets.Get("/", m.Handler.GetUserAccounts)
		wallets.Get("/:id", m.Handler.GetAccount)
		wallets.Post("/:id/deposit", m.Handler.Deposit)
		wallets.Post("/:id/withdraw", m.Handler.Withdraw)
		wallets.Post("/:id/transfer", m.Handler.Transfer)
		wallets.Get("/:id/transactions", m.Handler.GetTransactions)
	}
}
