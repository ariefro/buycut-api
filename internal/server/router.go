package server

import (
	"github.com/ariefro/buycut-api/internal/brand"
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/cronjobs"
	"github.com/ariefro/buycut-api/internal/middleware"
	"github.com/ariefro/buycut-api/internal/user"
	"github.com/gofiber/fiber/v2"
)

func setupRouter(
	app *fiber.App,
	userController user.Controller,
	companyController company.Controller,
	brandController brand.Controller,
) {
	api := app.Group("/api/v1")

	// users
	usersApi := api.Group("/users")
	usersApi.Post("/register", userController.Register)
	usersApi.Post("/login", userController.Login)

	// companies
	companiesApi := api.Group("/companies")
	companiesApi.Post("/", middleware.Auth(), companyController.Create)
	companiesApi.Get("/", companyController.Find)
	companiesApi.Put("/", middleware.Auth(), companyController.Update)
	companiesApi.Get("/:id", companyController.FindOneByID)
	companiesApi.Delete("/:id", middleware.Auth(), companyController.Delete)

	// brands
	brandsApi := api.Group("/brands")
	brandsApi.Post("/", middleware.Auth(), brandController.Create)
	brandsApi.Put("/:id", middleware.Auth(), brandController.Update)
	brandsApi.Delete("/:id", middleware.Auth(), brandController.Delete)

	brandsApi.Post("/boycotted", brandController.FindAll)
	brandsApi.Post("/search", brandController.FindByKeyword)

	// Cron Trigger
	cronjobs.Trigger(companyController)
}
