package product

import (
	"mime/multipart"

	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Create(c *fiber.Ctx) error
	FindByKeyword(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
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
	Name      string `form:"name" validate:"required~nama merek tidak boleh kosong"`
}

type createProductArgs struct {
	CompanyID  uint
	FormHeader *multipart.FileHeader
	Request    *createProductsRequest
}

type updateProductsRequest struct {
	Name      string `form:"name" validate:"required~nama merek tidak boleh kosong"`
	CompanyID *uint  `form:"company_id"`
}

type updateProductArgs struct {
	Product    *entity.Product
	Request    *updateProductsRequest
	FormHeader *multipart.FileHeader
}

type boycottedResult struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	Proof       []string        `json:"proof"`
	Company     *entity.Company `json:"company"`
	Type        string          `json:"type"` // Either "company" or "product"
}

type boycottedCountResult struct {
	CompanyCount int64 `json:"company_count"`
	ProductCount int64 `json:"product_count"`
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

	company, err := ctrl.companyService.FindOneByID(c.Context(), request.CompanyID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	formHeader, err := c.FormFile("image")
	if err != nil {
		response := helper.ResponseFailed("Image file is required")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctrl.service.Create(c.Context(), &createProductArgs{
		CompanyID:  company.ID,
		FormHeader: formHeader,
		Request:    &request,
	}); err != nil {
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

	res := helper.ResponseSuccess("Berhasil memuat daftar merek yang diboikot", result)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) FindAll(c *fiber.Ctx) error {
	var request getProductByKeywordRequest
	if err := c.BodyParser(&request); err != nil {
		res := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	count, err := ctrl.service.CountAll(c.Context(), &request)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	pages := pagination.NewFromRequest(c, int(count))
	paginationParams := pagination.PaginationParams{
		Offset: pages.Offset(),
		Limit:  pages.Size(),
	}

	results, err := ctrl.service.FindAll(c.Context(), &request, &paginationParams)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccessWithPagination("Berhasil memuat daftar merek yang diboikot", results, pages)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Update(c *fiber.Ctx) error {
	var request updateProductsRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := validator.ValidateStruct(request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	productID := helper.ParseStringToUint(c.Params("id"))
	product, err := ctrl.service.FindOneByID(c.Context(), productID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	formHeader, _ := c.FormFile("image")
	if err := ctrl.service.Update(c.Context(), productID, &updateProductArgs{
		Product:    product,
		Request:    &request,
		FormHeader: formHeader,
	}); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Data product berhasil diperbarui", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Delete(c *fiber.Ctx) error {
	productID := helper.ParseStringToUint(c.Params("id"))
	product, err := ctrl.service.FindOneByID(c.Context(), productID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	if err := ctrl.service.Delete(c.Context(), product); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menghapus merek dari daftar boikot", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}
