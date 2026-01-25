package utils

import (
	"errors"
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

	// Secret key untuk signing (harusnya dari environment variable)
	jwtSecret = []byte("your-secret-key-change-this-in-production")
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

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  15 * time.Minute,  // Access token: 15 menit
		refreshTokenDuration: 7 * 24 * time.Hour, // Refresh token: 7 hari
	}
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

func ExtractUserID(tokenString string) (uuid.UUID, error) {
	// Use default secret for backward compatibility
	manager := NewJWTManager(string(jwtSecret))
	claims, err := manager.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func SetSecretKey(secret string) {
	jwtSecret = []byte(secret)
}

func GetSecretKey() string {
	return string(jwtSecret)
}
