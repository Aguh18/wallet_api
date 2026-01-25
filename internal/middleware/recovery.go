package middleware

import (
	"wallet_api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// Recovery recovers from panics
func Recovery(l logger.Interface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				l.Error("Panic recovered: %v", r)
				c.Status(500).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    500,
						"message": "Internal server error",
					},
				})
			}
		}()
		return c.Next()
	}
}
