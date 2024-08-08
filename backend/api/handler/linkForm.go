package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thegroobi/web-listing-scrapper/models"
)

func LinkFormHandler(c *fiber.Ctx) error {
	c.Accepts("multipart/form-data", "application/json")
	u := new(models.UserData)

	err := c.BodyParser(u)
	if err != nil {
		log.Println(err.Error())
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "ok": false, "error": "Something went wrong"})
	}

	if !strings.HasPrefix(u.Link, "https://www.otomoto.pl/") {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "ok": false, "error": "Link is not an otomoto link"})
	}

	// TODO
	// save link with cookies since im the only one using it currently
	a := fiber.AcquireAgent()
	a.Cookies("link", u.Link)
	return c.Status(201).JSON(fiber.Map{"statusCode": 201, "ok": true, "message": "Link posted correctly"})
}
