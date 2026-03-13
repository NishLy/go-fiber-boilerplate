package auth

import (
	"go-boiler-plate/internal/user"
	"go-boiler-plate/pkg/utils"
	"go-boiler-plate/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService *JWTService
	DB          *gorm.DB
}

func NewAuthHandler(DB *gorm.DB, authService *JWTService) AuthHandler {
	return &authHandler{authService: authService, DB: DB}
}

func (a *authHandler) Login(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	var token string
	err := a.DB.Transaction(func(tx *gorm.DB) error {
		var user user.User
		if err := tx.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return err
		}
		// Here you should verify the password, for example using bcrypt
		// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 	return err
		// }
		t, err := a.authService.GenerateToken(user.ID)
		if err != nil {
			return err
		}
		token = t
		return nil
	})

	if err != nil {
		return utils.Unauthorized(c, "Invalid email or password")
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
