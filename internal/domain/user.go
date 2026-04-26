package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	TimestampModel
	SoftDeleteModel
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name"`
	Email           string    `json:"email" gorm:"unique"`
	Verified        bool      `json:"verified" gorm:"default:false"`
	Password        string    `json:"-"` // Exclude password from JSON responses
	Token           []Token   `gorm:"foreignKey:user_id;references:id" json:"-"`
	ProfileImageURL *string   `gorm:"type:text;nullable" json:"profile_image_url,omitempty"`
}

func (user *User) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
