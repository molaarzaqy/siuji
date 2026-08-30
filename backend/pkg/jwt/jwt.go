package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	PurposeVerifyEmail   = "verify_email"
	PurposeResetPassword = "reset_password"
)

type AccessTokenClaims struct {
	UserID   uint      `json:"user_id"`
	PublicID uuid.UUID `json:"pub_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID      uint   `json:"user_id"`
	TokenFamily string `json:"token_family"`
	jwt.RegisteredClaims
}

type TempTokenClaims struct {
	UserID  uint   `json:"user_id"`
	Purpose string `json:"purpose"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

// Manager holds JWT configuration and issues/validates tokens.
// Config is injected once at startup instead of read from a global variable.
type Manager struct {
	secret               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewManager(secret string, accessTokenDuration, refreshTokenDuration time.Duration) *Manager {
	return &Manager{
		secret:               secret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *Manager) AccessTokenDuration() time.Duration {
	return m.accessTokenDuration
}

func (m *Manager) RefreshTokenDuration() time.Duration {
	return m.refreshTokenDuration
}

func (m *Manager) GenerateAccessToken(userID uint, publicID uuid.UUID, email, role string) (string, error) {
	if m.secret == "" {
		return "", errors.New("JWT secret is required")
	}
	now := time.Now()
	claims := AccessTokenClaims{
		UserID:   userID,
		PublicID: publicID,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "siuji-backend",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *Manager) GenerateRefreshToken(userID uint, tokenFamily string) (string, error) {
	if m.secret == "" {
		return "", errors.New("JWT secret is required")
	}
	now := time.Now()
	claims := RefreshTokenClaims{
		UserID:      userID,
		TokenFamily: tokenFamily,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "siuji-backend",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *Manager) GenerateTempToken(userID uint, email, purpose string, expiryMinutes int) (string, error) {
	if m.secret == "" {
		return "", errors.New("JWT secret is required")
	}
	validPurposes := map[string]bool{
		PurposeVerifyEmail:   true,
		PurposeResetPassword: true,
	}
	if !validPurposes[purpose] {
		return "", fmt.Errorf("invalid purpose: %s", purpose)
	}
	now := time.Now()
	claims := TempTokenClaims{
		UserID:  userID,
		Email:   email,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "siuji-backend",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *Manager) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func (m *Manager) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func (m *Manager) ValidateTempToken(tokenString string, expectedPurpose string) (*TempTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TempTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*TempTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	if claims.Purpose != expectedPurpose {
		return nil, fmt.Errorf("invalid token purpose: expected %s, got %s", expectedPurpose, claims.Purpose)
	}
	return claims, nil
}

func GenerateTokenFamily() string {
	return uuid.New().String()
}