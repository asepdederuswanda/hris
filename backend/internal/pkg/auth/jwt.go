// Package auth menyediakan JWT authentication & authorization.
//
// Token structure:
//   - Access Token: short-lived (15 menit), bearer token untuk API calls
//   - Refresh Token: long-lived (24 jam), untuk mendapatkan access token baru
//
// JWT Claims mencakup user_id, company_id, role, dan permissions
// untuk multi-tenant context.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims adalah custom JWT claims untuk HRIS platform.
type Claims struct {
	UserID      string   `json:"sub"`
	CompanyID   string   `json:"company_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"` // "access" atau "refresh"
	jwt.RegisteredClaims
}

// Config untuk JWT.
type Config struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// Manager mengelola pembuatan dan validasi JWT token.
type Manager struct {
	cfg Config
}

// NewManager membuat JWT Manager baru.
func NewManager(cfg Config) *Manager {
	return &Manager{cfg: cfg}
}

// GenerateTokenPair membuat access dan refresh token untuk user.
func (m *Manager) GenerateTokenPair(userID, companyID, email, name, role string, permissions []string) (accessToken string, refreshToken string, err error) {
	// Access token (short-lived)
	accessToken, err = m.generateToken(userID, companyID, email, name, role, permissions, "access", m.cfg.AccessTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh token (long-lived)
	refreshToken, err = m.generateToken(userID, companyID, email, name, role, permissions, "refresh", m.cfg.RefreshTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken memvalidasi token dan mengembalikan claims.
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// RefreshToken memvalidasi refresh token dan menghasilkan access token baru.
func (m *Manager) RefreshToken(refreshTokenString string) (string, error) {
	claims, err := m.ValidateToken(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return "", fmt.Errorf("token is not a refresh token")
	}

	// Generate new access token
	accessToken, err := m.generateToken(
		claims.UserID, claims.CompanyID, claims.Email, claims.Name,
		claims.Role, claims.Permissions, "access", m.cfg.AccessTokenTTL,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return accessToken, nil
}

func (m *Manager) generateToken(userID, companyID, email, name, role string, permissions []string, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:      userID,
		CompanyID:   companyID,
		Email:       email,
		Name:        name,
		Role:        role,
		Permissions: permissions,
		TokenType:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.cfg.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{"hris-platform"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(m.cfg.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
