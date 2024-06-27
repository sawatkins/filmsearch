package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

type Movie struct {
	Title    string   `json:"title"`
	Year	 int	  `json:"year"`
}

type Movies struct {
	Movies []Movie `json:"movies"`
}

func Search(openaiClient *openai.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")
		movies, err := openaiMovieCompletion(openaiClient, query)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500) // maybe this should just dirent to a "try again later page"
		}

		return c.Render("search", fiber.Map{
			"Query":   query,
			"Results": unmarshallMovieTitles(movies),
		}, "layouts/main")
	}
}

func unmarshallMovieTitles(data string) []string {
	var movies Movies
	err := json.Unmarshal([]byte(data), &movies)
	if err != nil {
		fmt.Println(err)
		return nil // is this the best way to handle errors?
	}

	var movieTitles []string
	for _, movie := range movies.Movies {
		movieTitles = append(movieTitles, fmt.Sprintf("%s (%d)", movie.Title, movie.Year))
	}
	fmt.Println(movieTitles)
	return movieTitles
}
