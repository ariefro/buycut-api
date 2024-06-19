package company

import (
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Create(c *fiber.Ctx) error
	Find(c *fiber.Ctx) error
	FindOneByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service}
}

type createCompaniesRequest struct {
	Name        string `form:"name" validate:"required~nama perusahaan tidak boleh kosong"`
	Description string `form:"description"`
}

type getCompaniesRequest struct {
	Keyword string `json:"keyword"`
}

type updateCompaniesRequest struct {
	CompanyID uint   `json:"company_id" validate:"required~id tidak boleh kosong"`
	Name      string `json:"name" validate:"required~nama perusahaan tidak boleh kosong"`
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
	var request createCompaniesRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formHeader, _ := c.FormFile("image")
	if err := ctrl.service.Create(c.Context(), &request, formHeader); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menambahkan produk ke daftar boikot", nil)
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ctrl *controller) Find(c *fiber.Ctx) error {
	var request getCompaniesRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	count, err := ctrl.service.Count(c.Context(), &request)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	pages := pagination.NewFromRequest(c, int(count))
	paginationParams := pagination.PaginationParams{
		Offset: pages.Offset(),
		Limit:  pages.Size(),
	}

	result, err := ctrl.service.Find(c.Context(), &request, &paginationParams)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	data := helper.ResponseSuccessWithPagination("Berhasil memuat daftar produk yang diboikot", result, pages)
	return c.Status(fiber.StatusOK).JSON(data)
}

func (ctrl *controller) FindOneByID(c *fiber.Ctx) error {
	companyID := helper.ParseStringToUint(c.Params("id"))

	company, err := ctrl.service.FindOneByID(c.Context(), companyID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Produk ini masuk dalam daftar boikot!", company)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Update(c *fiber.Ctx) error {
	var request updateCompaniesRequest
	if err := c.BodyParser(&request); err != nil {
		res := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctrl.service.Update(c.Context(), &request); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil mengupdate data produk", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Delete(c *fiber.Ctx) error {
	companyID := helper.ParseStringToUint(c.Params("id"))

	if err := ctrl.service.Delete(c.Context(), companyID); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menghapus data produk", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}
