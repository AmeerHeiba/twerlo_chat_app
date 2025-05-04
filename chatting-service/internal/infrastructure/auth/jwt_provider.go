package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/config"
	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	cfg config.AuthConfig
}

func NewJWTProvider(cfg config.AuthConfig) *JWTProvider {
	return &JWTProvider{cfg: cfg}
}

func (p *JWTProvider) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	return p.generateToken(user, false)
}

func (p *JWTProvider) GenerateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	return p.generateToken(user, true)
}

func (p *JWTProvider) generateToken(user *domain.User, isRefresh bool) (string, error) {
	// Generate random session ID for invalidation purposes
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	expiry := p.cfg.AccessTokenExpiry
	if isRefresh {
		expiry = p.cfg.RefreshTokenExpiry
	}

	claims := domain.TokenClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		SessionID: sessionID,
		IsRefresh: isRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chatting-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.cfg.JWTSecret))
}

func (p *JWTProvider) ValidateToken(ctx context.Context, tokenString string) (*domain.TokenClaims, error) {
	claims, err := p.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.IsRefresh {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (p *JWTProvider) ValidateRefreshToken(ctx context.Context, tokenString string) (*domain.TokenClaims, error) {
	claims, err := p.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !claims.IsRefresh {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (p *JWTProvider) parseToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(p.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}

func (p *JWTProvider) GetAccessExpiry() time.Duration {
	return p.cfg.AccessTokenExpiry
}

func (p *JWTProvider) GetRefreshExpiry() time.Duration {
	return p.cfg.RefreshTokenExpiry
}

// generateSessionID creates a cryptographically secure random string
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
