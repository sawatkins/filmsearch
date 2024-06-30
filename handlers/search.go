package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

type Movie struct {
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Reason string `json:"justification"`
	// Poster string `json:"poster"`
}

type Movies struct {
	Movies []Movie `json:"movies"`
}

func Search(openaiClient *openai.Client, tmdbClient *tmdb.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")
		moviesJson, err := openaiMovieCompletion(openaiClient, query)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500) // maybe this should just dirent to a "try again later page"
		}

		movieTitles, movieReasons := unmarshallMovieTitles(moviesJson)

		// fmt.Println(movies.Movies[0].)

		posters := getTmdbInfo(tmdbClient, movieTitles)

		return c.Render("search", fiber.Map{
			"Query":   query,
			"Titles":  movieTitles,
			"Posters": posters,
			"Reasons": movieReasons,
			// "Movies": movies.Movies,
		}, "layouts/main")
	}
}

func getTmdbInfo(tmdbClient *tmdb.Client, movieTitles []string) []string {
	var posters []string
	for _, movieTitle := range movieTitles {
		parts := strings.Split(movieTitle, " (")
		if len(parts) < 2 {
			log.Println("Invalid movie title format:", movieTitle)
			continue
		}
		title := parts[0]
		year := strings.TrimSuffix(parts[1], ")")
		searchMovie, err := tmdbClient.GetSearchMovies(title, map[string]string{
			"primary_release_year": year,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		// TODO: make sure that this doesn't get added, or a blank things get added if not movie match found
		if len(searchMovie.Results) > 0 {
			posters = append(posters, tmdb.GetImageURL(searchMovie.Results[0].PosterPath, tmdb.W92))
		} else {
			log.Println("No results found for movie:", movieTitle)
		}
	}
	fmt.Println("posters:", posters)
	return posters
}

func unmarshallMovieTitles(data string) ([]string, []string) {
	var movies Movies
	err := json.Unmarshal([]byte(data), &movies)
	if err != nil {
		fmt.Println(err)
		return nil, nil // is this the best way to handle errors?
	}

	var movieTitles []string
	for _, movie := range movies.Movies {
		movieTitles = append(movieTitles, fmt.Sprintf("%s (%d)", movie.Title, movie.Year))
	}
	fmt.Println(movieTitles)

	var movieReasons []string
	for _, movie := range movies.Movies {
		movieReasons = append(movieReasons, movie.Reason)
	}
	fmt.Println(movieReasons)
	return movieTitles, movieReasons
}
