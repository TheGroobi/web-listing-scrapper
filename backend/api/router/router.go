package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thegroobi/web-listing-scrapper/api/handler"
)

func SetupRouter(app *fiber.App) {
	log.Println("Setting up routes")

	api := app.Group("/api")

	SetupCarListingsRoutes(api)
}

func SetupCarListingsRoutes(router fiber.Router) {
	cars := router.Group("/car-listings")
	cars.Get("/otomoto", handler.GetListings)
	cars.Patch("/link", handler.LinkFormHandler)
}
