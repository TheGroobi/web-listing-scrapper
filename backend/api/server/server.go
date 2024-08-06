package server

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartServer(port string) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server running!")
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
