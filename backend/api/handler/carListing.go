package handler

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/thegroobi/web-listing-scrapper/database"
	"github.com/thegroobi/web-listing-scrapper/models"
)

func GetListings(c *fiber.Ctx) error {
	db := db.DB
	var listings []models.CarListing

	db.Find(&listings)

	if len(listings) == 0 {
		return c.Status(404).JSON(fiber.Map{"statusCode": 404, "ok": false, "error": "No listings found"})
	}

	return c.Status(200).JSON(fiber.Map{"statusCode": 200, "ok": true, "message": "Listings found", "data": listings})
}
