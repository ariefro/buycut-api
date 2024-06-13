package product

import (
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Create(c *fiber.Ctx) error
}

type controller struct {
	service        Service
	companyService company.Service
}

func NewController(service Service, companyService company.Service) Controller {
	return &controller{service, companyService}
}

type createProductsRequest struct {
	CompanyID    uint     `json:"company_id" validate:"required~company id tidak boleh kosong"`
	ProductNames []string `json:"product_names" validate:"required~nama produk tidak boleh kosong"`
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
	var request createProductsRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := validator.ValidateStruct(request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	_, err := ctrl.companyService.FindOneByID(c.Context(), request.CompanyID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	if err := ctrl.service.Create(c.Context(), &request); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("berhasil menambahkan produk", nil)
	return c.Status(fiber.StatusCreated).JSON(res)
}
