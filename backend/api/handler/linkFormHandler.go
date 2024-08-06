package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/thegroobi/web-listing-scrapper/models"
)

func LinkFormHandler(c *fiber.Ctx) error {
	c.Accepts("multipart/form-data")
	u := new(models.UserData)

	err := c.BodyParser(u)
	if err != nil {
		log.Println(err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "Error", "code": 500, "message": "Something went wrong"})
	}

	a := fiber.AcquireAgent()

	// TODO
	// save link with cookies since im the only one using it currently
	a.Cookies("link", u.Link)
	return c.Status(201).JSON(fiber.Map{"status": "ok", "code": 201, "message": "Link posted correctly"})
}
