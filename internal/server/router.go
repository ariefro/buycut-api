package server

import (
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/gofiber/fiber/v2"
)

func setupRouter(
	app *fiber.App,
	companyController company.Controller,
) {
	api := app.Group("/api/v1")

	// company
	companiesApi := api.Group("/companies")
	companiesApi.Post("/", companyController.Create)
}
