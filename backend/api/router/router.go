package router

import (
	"github.com/gofiber/fiber/v2"
	listingHandler "github.com/web-listing-scrapper/api/handler"
)

func SetupRouter(app *fiber.App) {
	cars := app.Group("car_listings")
	cars.Get("/otomoto", listingHandler.GetListings)
}
