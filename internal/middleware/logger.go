package middleware

import (
	"time"

	"wallet_api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func Logger(l logger.Interface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		l.Info("[%s] %s %s - Status: %d - Duration: %v",
			c.Method(),
			c.Path(),
			c.IP(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
