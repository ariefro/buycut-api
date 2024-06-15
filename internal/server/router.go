package server

import (
	"github.com/ariefro/buycut-api/internal/company"
	"github.com/ariefro/buycut-api/internal/middleware"
	"github.com/ariefro/buycut-api/internal/product"
	"github.com/ariefro/buycut-api/internal/user"
	"github.com/gofiber/fiber/v2"
)

func setupRouter(
	app *fiber.App,
	userController user.Controller,
	companyController company.Controller,
	productController product.Controller,
) {
	api := app.Group("/api/v1")

	// users
	usersApi := api.Group("/users")
	usersApi.Post("/register", userController.Register)
	usersApi.Post("/login", userController.Login)

	// companies
	companiesApi := api.Group("/companies")
	companiesApi.Post("/", middleware.Auth(), companyController.Create)
	companiesApi.Get("/:id", companyController.FindOneByID)
	companiesApi.Put("/", companyController.Update)
	companiesApi.Delete("/:id", companyController.Delete)

	// products
	productsApi := api.Group("/products")
	productsApi.Post("/", middleware.Auth(), productController.Create)
	productsApi.Get("/boycotted-products", companyController.Find)
	productsApi.Post("/search", productController.FindByKeyword)
}
