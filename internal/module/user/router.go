package user

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers user routes
func (m *Module) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/v1/auth")
	{
		auth.Post("/register", m.Handler.Register)
		auth.Post("/login", m.Handler.Login)
	}

	users := app.Group("/v1/users")
	{
		users.Get("/profile", m.Handler.GetProfile)
		users.Put("/profile", m.Handler.UpdateProfile)
	}
}
