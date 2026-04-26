package user

import (
	"context"

	"github.com/NishLy/go-fiber-boilerplate/internal/domain"
	"go.uber.org/zap"
)

type UserService interface {
	GetUserFromEmail(ctx context.Context, email string) (*domain.User, error)
	RegisterUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
}

type userService struct {
	repo   UserRepository
	logger zap.SugaredLogger
}

func NewUserService(repo UserRepository, logger zap.SugaredLogger) UserService {
	return &userService{repo: repo, logger: logger}
}

func (s *userService) GetUserFromEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) RegisterUser(ctx context.Context, user *domain.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.DeleteUser(ctx, userID)
}
