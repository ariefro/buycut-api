package company

import (
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Create(c *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service}
}

type CreateCompaniesRequest struct {
	Names []string `json:"names" validate:"required~nama perusahaan tidak boleh kosong"`
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
	var createCompaniesRequest CreateCompaniesRequest
	if err := c.BodyParser(&createCompaniesRequest); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if errValid := validator.ValidateStruct(createCompaniesRequest); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctrl.service.Create(c.Context(), createCompaniesRequest); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("berhasil mendaftarkan perusahaan", nil)

	return c.Status(fiber.StatusCreated).JSON(res)
}
