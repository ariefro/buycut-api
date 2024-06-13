package server

import (
	"github.com/ariefro/buycut-api/config"
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/product"
	"github.com/ariefro/buycut-api/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	log "github.com/sirupsen/logrus"
)

func NewFiberServer(
	config *config.Config,
	userController user.Controller,
	companyController company.Controller,
	productController product.Controller,
) error {
	log.Println("starting server...")
	app := fiber.New()
	app.Use(recover.New())

	setupRouter(
		app,
		userController,
		companyController,
		productController,
	)

	log.Printf("🚀 listening on %s", config.AppPort)
	return app.Listen(":" + config.AppPort)
}
