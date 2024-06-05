package user

import (
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Register(c *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service}
}

type registerRequest struct {
	Name     string `json:"name" validate:"required~nama tidak boleh kosong"`
	Email    string `json:"email" validate:"required~email tidak boleh kosong"`
	Password string `json:"password" validate:"required~password tidak boleh kosong"`
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
