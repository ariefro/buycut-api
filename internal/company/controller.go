package company

import (
	"context"
	"mime/multipart"

	"github.com/ariefro/buycut-api/internal/entity"
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/ariefro/buycut-api/pkg/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/usepzaka/validator"
)

type Controller interface {
	Create(c *fiber.Ctx) error
	Find(c *fiber.Ctx) error
	FindOneByID(c *fiber.Ctx) error
	FindOneDummy(context.Context) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service}
}

type createCompanyRequest struct {
	Name        string   `form:"name" validate:"required~nama perusahaan tidak boleh kosong"`
	Description string   `form:"description" validate:"required~deskripsi tidak boleh kosong"`
	Proof       []string `form:"proof" validate:"required~bukti tidak boleh kosong"`
}

type createCompanyArgs struct {
	FormHeader *multipart.FileHeader
	Request    *createCompanyRequest
}

type getCompaniesRequest struct {
	Keyword string `json:"keyword"`
}

type updateCompanyRequest struct {
	CompanyID   uint     `form:"company_id" validate:"required~company id tidak boleh kosong"`
	Name        *string  `form:"name"`
	Description *string  `form:"description"`
	ImageURL    *string  `form:"image_url"`
	Proof       []string `form:"proof"`
}

type updateCompanyArgs struct {
	Company    *entity.Company
	FormHeader *multipart.FileHeader
	Request    *updateCompanyRequest
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
	var request createCompanyRequest
	if err := c.BodyParser(&request); err != nil {
		response := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	formHeader, err := c.FormFile("image")
	if err != nil {
		response := helper.ResponseFailed("Image file is required")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctrl.service.Create(c.Context(), &createCompanyArgs{
		FormHeader: formHeader,
		Request:    &request,
	}); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menambahkan merek ke daftar boikot", nil)
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ctrl *controller) Find(c *fiber.Ctx) error {
	count, err := ctrl.service.Count(c.Context())
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	pages := pagination.NewFromRequest(c, int(count))
	paginationParams := pagination.PaginationParams{
		Offset: pages.Offset(),
		Limit:  pages.Size(),
	}

	result, err := ctrl.service.Find(c.Context(), &paginationParams)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	data := helper.ResponseSuccessWithPagination("Berhasil memuat daftar merek yang diboikot", result, pages)
	return c.Status(fiber.StatusOK).JSON(data)
}

func (ctrl *controller) FindOneByID(c *fiber.Ctx) error {
	companyID := helper.ParseStringToUint(c.Params("id"))

	company, err := ctrl.service.FindOneByID(c.Context(), companyID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Merek ini masuk dalam daftar boikot!", company)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) FindOneDummy(ctx context.Context) error {
	_, err := ctrl.service.FindOneByID(ctx, 1)
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *controller) Update(c *fiber.Ctx) error {
	var request updateCompanyRequest
	if err := c.BodyParser(&request); err != nil {
		res := helper.ResponseFailed(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	if errValid := validator.ValidateStruct(request); errValid != nil {
		response := helper.ResponseFailed(errValid.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	company, err := ctrl.service.FindOneByID(c.Context(), request.CompanyID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	formHeader, _ := c.FormFile("image")
	if err := ctrl.service.Update(c.Context(), &updateCompanyArgs{
		Company:    company,
		FormHeader: formHeader,
		Request:    &request,
	}); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil mengupdate data merek", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (ctrl *controller) Delete(c *fiber.Ctx) error {
	companyID := helper.ParseStringToUint(c.Params("id"))
	company, err := ctrl.service.FindOneByID(c.Context(), companyID)
	if err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	if err := ctrl.service.Delete(c.Context(), company); err != nil {
		return helper.GenerateErrorResponse(c, err.Error())
	}

	res := helper.ResponseSuccess("Berhasil menghapus merek dari daftar boikot", nil)
	return c.Status(fiber.StatusOK).JSON(res)
}
