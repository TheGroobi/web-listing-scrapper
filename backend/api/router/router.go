package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thegroobi/web-listing-scrapper/api/handler"
)

func SetupRouter(app *fiber.App) {
	log.Println("Setting up routes")

	api := app.Group("/api")

	SetupCarListingsRoutes(api)
}

func SetHeaders(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
}

func SetupCarListingsRoutes(router fiber.Router) {
	cars := router.Group("/car-listings")
	cars.Get("/otomoto", handler.GetListings)
	cars.Post("/link", handler.LinkFormHandler)
}
