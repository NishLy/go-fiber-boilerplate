package user

import (
	"context"

	"github.com/NishLy/go-fiber-boilerplate/internal/database"
	"github.com/NishLy/go-fiber-boilerplate/internal/domain"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type userRepository struct {
	logger zap.SugaredLogger
}

func NewUserRepository(log zap.SugaredLogger) UserRepository {
	return &userRepository{logger: log}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	db, err := database.GetDB(database.GetIndentifier(ctx), false)
	if err != nil {
		return database.Wrap(err)
	}

	err = db.DB.Create(user).Error

	if err != nil {
		r.logger.Errorf("Failed to create user: %v", err)
		return database.Wrap(err)
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db, err := database.GetDB(database.GetIndentifier(ctx), false)
	if err != nil {
		return nil, database.Wrap(err)
	}

	var user domain.User
	err = db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		r.logger.Errorf("Failed to get user by email: %v", err)
		return nil, database.Wrap(err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	db, err := database.GetDB(database.GetIndentifier(ctx), false)
	if err != nil {
		return nil, database.Wrap(err)
	}

	var user domain.User
	err = db.DB.First(&user, id).Error
	if err != nil {
		r.logger.Errorf("Failed to get user by ID: %v", err)
		return nil, database.Wrap(err)
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	db, err := database.GetDB(database.GetIndentifier(ctx), false)
	if err != nil {
		return database.Wrap(err)
	}

	err = db.DB.Save(user).Error
	if err != nil {
		r.logger.Errorf("Failed to update user: %v", err)
		return database.Wrap(err)
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	db, err := database.GetDB(database.GetIndentifier(ctx), false)
	if err != nil {
		return database.Wrap(err)
	}
	err = db.DB.Delete(&domain.User{}, id).Error
	if err != nil {
		r.logger.Errorf("Failed to delete user: %v", err)
		return database.Wrap(err)
	}
	return nil
}
