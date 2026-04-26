package auth

import (
	"github.com/NishLy/go-fiber-boilerplate/pkg/utils"
	"github.com/NishLy/go-fiber-boilerplate/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService *authService
}

func NewAuthHandler(authService *authService) AuthHandler {
	return &authHandler{authService: authService}
}

func (a *authHandler) Login(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	token := "token" // Placeholder for token generation logic

	return c.JSON(fiber.Map{
		"token": token,
	})
}
