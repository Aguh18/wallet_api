package user

import (
	"wallet_api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/v1/auth")
	{
		auth.Post("/register", m.Handler.Register)
		auth.Post("/login", m.Handler.Login)
		auth.Post("/logout", middleware.JWTAuth(), m.Handler.Logout)
		auth.Post("/refresh", m.Handler.RefreshToken)
	}

	users := app.Group("/v1/users", middleware.JWTAuth())
	{
		users.Get("/profile", m.Handler.GetProfile)
		users.Put("/profile", m.Handler.UpdateProfile)
	}
}
