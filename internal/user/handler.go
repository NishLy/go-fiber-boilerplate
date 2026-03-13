package user

import "gorm.io/gorm"

type UserHandler interface {
	GetUsers() ([]User, error)
}

type userHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return &userHandler{DB: db}
}

// GetUsers implements [UserHandler].
func (u *userHandler) GetUsers() ([]User, error) {
	panic("unimplemented")
}