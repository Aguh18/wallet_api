package utils

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// ErrInvalidToken menandakan token tidak valid
	ErrInvalidToken = errors.New("token is invalid")

	// ErrExpiredToken menandakan token sudah expired
	ErrExpiredToken = errors.New("token has expired")

	// ErrTokenMalformed menandakan token format salah
	ErrTokenMalformed = errors.New("token is malformed")
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // access token expiry in seconds
}

type JWTManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTManager creates JWT manager with configuration from environment variables
func NewJWTManager(secretKey string) *JWTManager {
	// Baca konfigurasi dari environment variables
	accessTokenExpiry := getEnvAsInt("ACCESS_TOKEN_EXPIRY", 15)  // default 15 menit
	refreshTokenExpiry := getEnvAsInt("REFRESH_TOKEN_EXPIRY", 7)  // default 7 hari

	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  time.Duration(accessTokenExpiry) * time.Minute,
		refreshTokenDuration: time.Duration(refreshTokenExpiry) * 24 * time.Hour,
	}
}

// getEnvAsInt reads environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func (j *JWTManager) GenerateToken(userID uuid.UUID, username string) (*TokenPair, error) {
	// Generate Access Token
	accessTokenExpiry := time.Now().Add(j.accessTokenDuration)
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wallet_api",
			Subject:   userID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshTokenExpiry := time.Now().Add(j.refreshTokenDuration)
	refreshClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wallet_api",
			Subject:   userID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(j.accessTokenDuration.Seconds()),
	}, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrInvalidToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (j *JWTManager) RefreshAccessToken(refreshTokenString string) (string, error) {
	// Validate refresh token
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Generate new access token dengan claims yang sama
	newAccessTokenExpiry := time.Now().Add(j.accessTokenDuration)
	newAccessClaims := &Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(newAccessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wallet_api",
			Subject:   claims.UserID.String(),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	return newAccessToken.SignedString(j.secretKey)
}

// GetSecretKey returns JWT secret from environment variable
func GetSecretKey() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	// Fallback ke default untuk development dengan WARNING
	log.Println("⚠️  WARNING: Menggunakan default JWT secret! Set JWT_SECRET environment variable untuk production!")
	return "your-secret-key-change-this-in-production"
}
