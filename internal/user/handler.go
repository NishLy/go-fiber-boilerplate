package user

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/domain"
	"gorm.io/gorm"
)

type UserHandler interface {
	GetUsers() ([]domain.User, error)
}

type userHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return &userHandler{DB: db}
}

// GetUsers implements [UserHandler].
func (u *userHandler) GetUsers() ([]domain.User, error) {
	panic("unimplemented")
}
