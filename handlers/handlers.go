package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, Index!",
	}, "layouts/main")
}

func About(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Hello, About!",
	}, "layouts/main")
}

