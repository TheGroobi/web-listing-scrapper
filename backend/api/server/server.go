package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartServer() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server running!")
	})

	log.Fatal(app.Listen(":3000"))
}
