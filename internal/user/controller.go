package user

import (
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service}
}

type registerRequest struct {
	Name     string `json:"name" validate:"required~nama tidak boleh kosong"`
	Email    string `json:"email" validate:"required~email tidak boleh kosong, email~format email tidak valid"`
	Password string `json:"password" validate:"required~password tidak boleh kosong"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required~email tidak boleh kosong, email~format email tidak valid"`
	Password string `json:"password" validate:"required~password tidak boleh kosong"`
}

type loginResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func (ctrl *controller) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if errValid := validator.ValidateStruct(req); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctrl.service.Register(c.Context(), &req); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("registrasi user berhasil", nil)
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ctrl *controller) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if errValid := validator.ValidateStruct(req); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	user, err := ctrl.service.FindOneByEmail(c.Context(), req.Email)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	accessToken, err := ctrl.service.Login(c.Context(), &req, user)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	data := &loginResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}

	response := helper.ResponseSuccess("login berhasil", data)
	return c.Status(fiber.StatusOK).JSON(response)
}
