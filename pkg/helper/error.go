package helper

import (
	"github.com/ariefro/buycut-api/pkg/common"
	"github.com/gofiber/fiber/v2"
)

func GenerateErrorResponse(c *fiber.Ctx, errorMessage string) error {
	var statusCode int

	switch errorMessage {
	case common.ErrDuplicateEntry:
		statusCode = fiber.StatusConflict
	default:
		statusCode = fiber.StatusInternalServerError
	}

	resp := ResponseFailed(errorMessage)
	return c.Status(statusCode).JSON(resp)
}
