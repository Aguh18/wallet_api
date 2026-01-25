package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	// Cookie names
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

type CookieConfig struct {
	Name     string
	Value    string
	MaxAge   time.Duration
	HTTPOnly bool
	Secure   bool
	SameSite string
	Path     string
}

func SetCookie(c *fiber.Ctx, config CookieConfig) {
	// Default values
	if config.Path == "" {
		config.Path = "/"
	}
	if config.SameSite == "" {
		config.SameSite = "Strict"
	}

	cookie := &fiber.Cookie{
		Name:     config.Name,
		Value:    config.Value,
		MaxAge:   int(config.MaxAge.Seconds()),
		HTTPOnly: config.HTTPOnly,
		Secure:   config.Secure,
		SameSite: config.SameSite, // "Strict", "Lax", or "None"
		Path:     config.Path,
	}

	c.Cookie(cookie)
}

func SetAuthCookies(c *fiber.Ctx, accessToken, refreshToken string, accessTokenExpiry time.Duration) {
	// Access Token cookie (httpOnly, secure)
	SetCookie(c, CookieConfig{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		MaxAge:   accessTokenExpiry,
		HTTPOnly: true,
		Secure:   true,  // Only HTTPS
		SameSite: "Strict",
		Path:     "/",
	})

	// Refresh Token cookie (httpOnly, secure, longer expiry)
	SetCookie(c, CookieConfig{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		MaxAge:   7 * 24 * time.Hour, // 7 days
		HTTPOnly: true,
		Secure:   true,  // Only HTTPS
		SameSite: "Strict",
		Path:     "/",
	})
}

func ClearAuthCookies(c *fiber.Ctx) {
	// Clear access token cookie
	SetCookie(c, CookieConfig{
		Name:     AccessTokenCookie,
		Value:    "",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	})

	// Clear refresh token cookie
	SetCookie(c, CookieConfig{
		Name:     RefreshTokenCookie,
		Value:    "",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	})
}

func GetAccessTokenFromCookie(c *fiber.Ctx) string {
	return c.Cookies(AccessTokenCookie)
}

func GetRefreshTokenFromCookie(c *fiber.Ctx) string {
	return c.Cookies(RefreshTokenCookie)
}

func IsDevelopmentMode(app *fiber.App) bool {
	// You can also check environment variable
	// For now, assume production when not explicitly set
	return false // Can be configured via env var
}

func SetAuthCookiesDevelopment(c *fiber.Ctx, accessToken, refreshToken string, accessTokenExpiry time.Duration) {
	// Access Token cookie (development mode - tanpa secure flag)
	SetCookie(c, CookieConfig{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		MaxAge:   accessTokenExpiry,
		HTTPOnly: true,
		Secure:   false, // HTTP allowed in development
		SameSite: "Lax",
		Path:     "/",
	})

	// Refresh Token cookie (development mode)
	SetCookie(c, CookieConfig{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		MaxAge:   7 * 24 * time.Hour,
		HTTPOnly: true,
		Secure:   false, // HTTP allowed in development
		SameSite: "Lax",
		Path:     "/",
	})
}

func SetAuthCookiesSmart(c *fiber.Ctx, accessToken, refreshToken string, accessTokenExpiry time.Duration, isProduction bool) {
	if isProduction {
		SetAuthCookies(c, accessToken, refreshToken, accessTokenExpiry)
	} else {
		SetAuthCookiesDevelopment(c, accessToken, refreshToken, accessTokenExpiry)
	}
}
