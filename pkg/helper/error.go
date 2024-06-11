package helper

import (
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GenerateErrorResponse(c *fiber.Ctx, errorMessage string) error {
	var statusCode int

	switch errorMessage {
	case common.ErrInvalidEmailOrPassword:
		statusCode = fiber.StatusBadRequest
	case common.MissingJWT:
		statusCode = fiber.StatusUnauthorized
	case common.EmailNotRegistered,
		common.CompanyNotFound:
		statusCode = fiber.StatusNotFound
	case common.ErrDuplicateEntry,
		gorm.ErrDuplicatedKey.Error():
		statusCode = fiber.StatusConflict
		errorMessage = common.ErrDuplicateEntry
	default:
		statusCode = fiber.StatusInternalServerError
	}

	resp := ResponseFailed(errorMessage)
	return c.Status(statusCode).JSON(resp)
}
