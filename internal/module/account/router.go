package account

import (
	"wallet_api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(app *fiber.App) {
	accounts := app.Group("/v1/accounts", middleware.JWTAuth())
	{

		accounts.Post("/", m.Handler.CreateAccount)
		accounts.Get("/", m.Handler.GetUserAccounts)
		accounts.Get("/:id", m.Handler.GetAccount)
		accounts.Post("/:id/deposit", m.Handler.Deposit)
		accounts.Post("/:id/withdraw", m.Handler.Withdraw)
		accounts.Get("/:id/transactions", m.Handler.GetTransactions)
	}
}
