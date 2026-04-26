package token

import (
	"context"
	"time"

	"github.com/NishLy/go-fiber-boilerplate/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type TokenService interface {
	GenerateToken(ctx context.Context, userID string, expires time.Duration, tokenType string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID string) (string, error)
	GenerateAccessToken(ctx context.Context, userID string) (string, error)
	DeleteTokenByUserID(ctx context.Context, userID string, tokenType string) error
	GenerateForgotPasswordToken(ctx context.Context, userID string) (string, error)
}

type tokenService struct {
	logger zap.SugaredLogger
	r      TokenRepository
}

func NewTokenService(logger zap.SugaredLogger, r TokenRepository) TokenService {
	return &tokenService{
		logger: logger,
		r:      r,
	}
}

func (s *tokenService) GenerateToken(ctx context.Context, userID string, expires time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expires).Unix(),
		"type":    tokenType,
	}

	config := config.Get()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		s.logger.Errorf("Failed to sign token: %v", err)
		return "", err
	}

	err = s.r.SaveToken(ctx, userID, tokenStr, tokenType, expires)
	if err != nil {
		s.logger.Errorf("Failed to save token: %v", err)
		return "", err
	}

	return tokenStr, nil
}

func (s *tokenService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	return s.GenerateToken(ctx, userID, time.Hour*24*7, TokenTypeRefresh)
}

func (s *tokenService) GenerateAccessToken(ctx context.Context, userID string) (string, error) {
	return s.GenerateToken(ctx, userID, time.Hour*1, TokenTypeAccess)
}

func (s *tokenService) DeleteTokenByUserID(ctx context.Context, userID string, tokenType string) error {
	return s.r.DeleteTokenByUserID(ctx, userID, tokenType)
}

func (s *tokenService) GenerateForgotPasswordToken(ctx context.Context, userID string) (string, error) {
	return s.GenerateToken(ctx, userID, time.Hour*1, TokenTypeResetPassword)
}
