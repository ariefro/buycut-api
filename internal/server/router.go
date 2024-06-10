package server

import (
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/middleware"
	"github.com/ariefro/buycut-api/internal/user"
	"github.com/gofiber/fiber/v2"
)

func setupRouter(
	app *fiber.App,
	userController user.Controller,
	companyController company.Controller,
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
}
