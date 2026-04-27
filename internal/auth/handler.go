package auth

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/response"
	"github.com/NishLy/go-fiber-boilerplate/pkg/utils"
	"github.com/NishLy/go-fiber-boilerplate/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type authHandler struct {
	authService *authService
}

func NewAuthHandler(authService *authService) AuthHandler {
	return &authHandler{authService: authService}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (a *authHandler) Login(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	token := "token" // Placeholder for token generation logic

	return c.Status(fiber.StatusOK).
		JSON(response.GenericSuccessResponse[fiber.Map]{
			GenericResponse: response.GenericResponse{
				Message: "Login successful",
			},
			Data: fiber.Map{
				"token": token,
			},
		})
}

// Register godoc
// @Summary User registration
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (a *authHandler) Register(c *fiber.Ctx) error {

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.GenericResponse{
			Message: "Registration successful",
		})
}
