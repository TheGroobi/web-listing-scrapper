package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thegroobi/web-listing-scrapper/api/router"
	"github.com/thegroobi/web-listing-scrapper/config"
	db "github.com/thegroobi/web-listing-scrapper/database"
)

var link = "https://www.otomoto.pl/osobowe/honda/accord"

func main() {
	cfg := config.LoadConfig()
	app := fiber.New()

	db.InitDB(cfg)
	db.DBStatus()

	// otomoto.ScrapArticles(link)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server running!")
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.ServerPort)))
	router.SetupRouter(app)
}
