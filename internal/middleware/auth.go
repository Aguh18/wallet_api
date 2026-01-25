package middleware

import (
	"wallet_api/internal/common/response"
	"wallet_api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookie only
		tokenString := utils.GetAccessTokenFromCookie(c)

		// If no token, return unauthorized
		if tokenString == "" {
			return c.Status(401).JSON(response.Error(401, "Authentication required"))
		}

		// 4. Validate token
		jwtManager := utils.NewJWTManager(utils.GetSecretKey())
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			switch err {
			case utils.ErrExpiredToken:
				return c.Status(401).JSON(response.Error(401, "Token has expired"))
			case utils.ErrInvalidToken, utils.ErrTokenMalformed:
				return c.Status(401).JSON(response.Error(401, "Invalid token"))
			default:
				return c.Status(401).JSON(response.Error(401, "Authentication failed"))
			}
		}

		// 5. Store user info in context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

// Berguna untuk routes yang bisa diakses public tapi dengan extra features jika logged in
func OptionalJWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get token from cookie
		tokenString := utils.GetAccessTokenFromCookie(c)

		// If we have a token, validate it
		if tokenString != "" {
			jwtManager := utils.NewJWTManager(utils.GetSecretKey())
			claims, err := jwtManager.ValidateToken(tokenString)
			if err == nil {
				// Token valid, store user info
				c.Locals("user_id", claims.UserID)
				c.Locals("username", claims.Username)
				c.Locals("authenticated", true)
			} else {
				// Token invalid but we don't block the request
				c.Locals("authenticated", false)
			}
		} else {
			c.Locals("authenticated", false)
		}

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	return userID, ok
}

func GetUsername(c *fiber.Ctx) (string, bool) {
	username, ok := c.Locals("username").(string)
	return username, ok
}

func IsAuthenticated(c *fiber.Ctx) bool {
	authenticated, ok := c.Locals("authenticated").(bool)
	if !ok {
		// Fallback: check if user_id exists
		_, hasUserID := GetUserID(c)
		return hasUserID
	}
	return authenticated
}
