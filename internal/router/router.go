package router

import (
	"net/http"

	"wallet_api/config"
	"wallet_api/internal/middleware"
	"wallet_api/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	routerModule *Module,
	l logger.Interface,
) {
	// Global middleware
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(cors.New())
	app.Use(recover.New())

	// Health check endpoint (Kubernetes standard)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Register all module routes
	routerModule.RegisterRoutes(app)
}
