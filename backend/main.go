package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thegroobi/web-listing-scrapper/api/router"
	"github.com/thegroobi/web-listing-scrapper/config"
	db "github.com/thegroobi/web-listing-scrapper/database"
)

func main() {
	cfg := config.LoadConfig()
	app := fiber.New()

	db.InitDB(cfg)
	db.DBStatus()

	router.SetHeaders(app)
	router.SetupRouter(app)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.ServerPort)))
}
