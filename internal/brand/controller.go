package brand

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

type getBrandByKeywordRequest struct {
	Keyword string `json:"keyword"`
}

type createBrandsRequest struct {
	CompanyID uint   `form:"company_id" validate:"required~company id tidak boleh kosong"`
	Name      string `form:"name" validate:"required~nama merek tidak boleh kosong"`
}

type createBrandArgs struct {
	CompanyID  uint
	FormHeader *multipart.FileHeader
	Request    *createBrandsRequest
}

type updateBrandsRequest struct {
	Name      string `form:"name" validate:"required~nama merek tidak boleh kosong"`
	CompanyID *uint  `form:"company_id"`
}

type updateBrandArgs struct {
	Brand      *entity.Brand
	Request    *updateBrandsRequest
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
	Type        string          `json:"type"` // Either "company" or "brand"
}

type boycottedCountResult struct {
	CompanyCount int64 `json:"company_count"`
	BrandCount   int64 `json:"brand_count"`
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
	var request createBrandsRequest
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

	if err := ctrl.service.Create(c.Context(), &createBrandArgs{
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
	var request getBrandByKeywordRequest
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
	var request getBrandByKeywordRequest
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

	// pages limit for query companies and brands
	pages.Limit = pages.Limit * 2
	res := helper.ResponseSuccessWithPagination("Berhasil memuat daftar merek yang diboikot", results, pages)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Update(c *fiber.Ctx) error {
	var request updateBrandsRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := validator.ValidateStruct(request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	brandID := helper.ParseStringToUint(c.Params("id"))
	brand, err := ctrl.service.FindOneByID(c.Context(), brandID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	formHeader, _ := c.FormFile("image")
	if err := ctrl.service.Update(c.Context(), brandID, &updateBrandArgs{
		Brand:      brand,
		Request:    &request,
		FormHeader: formHeader,
	}); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Data brand berhasil diperbarui", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Delete(c *fiber.Ctx) error {
	brandID := helper.ParseStringToUint(c.Params("id"))
	brand, err := ctrl.service.FindOneByID(c.Context(), brandID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	if err := ctrl.service.Delete(c.Context(), brand); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menghapus merek dari daftar boikot", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}
