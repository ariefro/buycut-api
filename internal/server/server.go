package server

import (
	"github.com/ariefro/buycut-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	log "github.com/sirupsen/logrus"
)

func NewFiberServer(config *config.Config) error {
	log.Println("starting server...")
	app := fiber.New()
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Printf("🚀 listening on %s", config.AppPort)
	return app.Listen(":" + config.AppPort)
}
