package auth

import (
	"context"

	"github.com/NishLy/go-fiber-boilerplate/internal/domain"
	"github.com/NishLy/go-fiber-boilerplate/internal/token"
	"github.com/NishLy/go-fiber-boilerplate/internal/user"
	"github.com/NishLy/go-fiber-boilerplate/pkg"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, req RegisterRequest) (string, error)
}

type authService struct {
	userService  user.UserService
	tokenService token.TokenService
}

func NewAuthService(userService user.UserService, tokenService token.TokenService) *authService {
	return &authService{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (j *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := j.userService.GetUserFromEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !pkg.CheckPasswordHash(password, user.Password) {
		return "", err
	}

	tokenStr, err := j.tokenService.GenerateAccessToken(ctx, user.ID.String())
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *authService) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := j.userService.RegisterUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
