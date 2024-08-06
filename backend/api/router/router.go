package router

import (
	"github.com/gofiber/fiber/v2"
	carListingHandler "github.com/thegroobi/web-listing-scrapper/api/handler"
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/api")

	SetupCarListingsRoutes(api)
}

func SetupCarListingsRoutes(router fiber.Router) {
	cars := router.Group("/car-listings")
	cars.Get("/otomoto", carListingHandler.GetListings)
}
