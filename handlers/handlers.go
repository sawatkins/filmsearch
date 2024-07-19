package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).Render("404", fiber.Map{
		"Message": "404 Not found! Please try again",
	}, "layouts/main")
}

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":       "FilmSearch",
		"Canonical":   "https://filmsearch.xyz",
		"Robots":      "index, follow",
		"Description": "AI search engine to discover movies using natural language",
		"Keywords":    "filmsearch, film, search, movie, discover, ai",
	}, "layouts/main")
}

func About(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title":       "FilmSearch - About",
		"Canonical":   "https://filmsearch.xyz/about",
		"Robots":      "index, follow",
		"Description": "About FilmSearch",
		"Keywords":    "filmsearch, about, film, search, movie, discover, ai",
	}, "layouts/main")
}
