package middleware

import (
	"github.com/ariefro/buycut-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/spf13/viper"
)

func Auth() fiber.Handler {
	config := jwtware.Config{
		SigningKey:   []byte(viper.GetString("JWT_SECRET_KEY")),
		ErrorHandler: jwtError,
	}

	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	return helper.GenerateErrorResponse(c, err.Error())
}
