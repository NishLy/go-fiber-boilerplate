package user

import (
	"context"
	"fmt"

	"github.com/NishLy/go-fiber-boilerplate/internal/domain"
	"github.com/NishLy/go-fiber-boilerplate/internal/platform/database"
	"github.com/NishLy/go-fiber-boilerplate/internal/request"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context, pagination request.PaginationRequest) ([]domain.User, paginator.Cursor, error)
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

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
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

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
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

func (r *userRepository) GetUsers(ctx context.Context, pagination request.PaginationRequest) ([]domain.User, paginator.Cursor, error) {
	var users []domain.User

	fmt.Printf("Get db from context for GetUsers: %v %s\n ", ctx.Value("db"), database.GetIndentifier(ctx))

	db, err := database.GetDBFromContext(ctx)
	if err != nil {
		return nil, paginator.Cursor{}, err
	}

	p := paginator.New(&paginator.Config{
		// clean up the sort_by input to prevent SQL injection
		Keys:  []string{domain.GetSortColumn(pagination.SortBy)},
		Limit: pagination.Limit,
		Order: paginator.Order(pagination.Sort),
	})

	p.SetAfterCursor(pagination.AfterCursor)

	query := db.Model(&domain.User{})

	if pagination.Search != "" {
		searchTerm := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	result, cursor, err := p.Paginate(query, &users)
	if err != nil {
		return nil, paginator.Cursor{}, database.Wrap(err)
	}

	if result.Error != nil {
		return nil, paginator.Cursor{}, database.Wrap(result.Error)
	}

	return users, cursor, nil
}
