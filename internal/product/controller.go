package product

import (
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Create(c *fiber.Ctx) error
	FindByKeyword(c *fiber.Ctx) error
}

type controller struct {
	service        Service
	companyService company.Service
}

func NewController(service Service, companyService company.Service) Controller {
	return &controller{service, companyService}
}

type getProductByKeywordRequest struct {
	Keyword string `json:"keyword"`
}

type createProductsRequest struct {
	CompanyID uint   `form:"company_id" validate:"required~company id tidak boleh kosong"`
	Name      string `form:"name" validate:"required~nama perusahaan tidak boleh kosong"`
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

	formHeader, err := c.FormFile("image")
	if err != nil {
		response := helper.ResponseFailed("Image file is required")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctrl.service.Create(c.Context(), &request, formHeader); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menambahkan merek ke daftar boikot", nil)
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ctrl *controller) FindByKeyword(c *fiber.Ctx) error {
	var request getProductByKeywordRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	result, err := ctrl.service.FindByKeyword(c.Context(), &request)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Merek ini masuk dalam daftar boikot!", result)
	return c.Status(fiber.StatusOK).JSON(res)
}
