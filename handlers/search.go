package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	tmdb "github.com/cyruzin/golang-tmdb"
	openai "github.com/sashabaranov/go-openai"
)

type Movie struct {
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Reason string `json:"justification"`
}

type Movies struct {
	Movies []Movie `json:"movies"`
}

func Search(openaiClient *openai.Client, tmdbClient *tmdb.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")

		if query == "" {
			return c.Render("404", fiber.Map{
				"Message": "No query provided",
			}, "layouts/main")
		}

		return c.Render("search", fiber.Map{
			"Query":   query,
		}, "layouts/main")
	}
}

func SearchResults(openaiClient *openai.Client, tmdbClient *tmdb.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")
		moviesJson, err := openaiMovieCompletion(openaiClient, query)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500) // maybe this should just dirent to a "try again later page"
		}

		movieTitles, movieReasons := unmarshallMovieTitles(moviesJson)
		posters, tmdbUrls := getTmdbInfo(tmdbClient, movieTitles)

		return c.Render("search-results", fiber.Map{
			"Titles":  movieTitles,
			"Posters": posters,
			"Reasons": movieReasons,
			"Urls":    tmdbUrls,
		})
	}
}

func getTmdbInfo(tmdbClient *tmdb.Client, movieTitles []string) ([]string, []string) {
	var posters []string
	var tmdbUrls []string
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
			log.Println("Error searching for movie:", err)
			continue
		}

		if len(searchMovie.Results) == 0 {
			log.Println("No results found for movie:", movieTitle)
			posters = append(posters, "/img/no_movie_poster_found.jpg")
			tmdbUrls = append(tmdbUrls, "notfound")
			continue
		}

		// get url
		baseURL := "https://www.themoviedb.org/movie/"
		tmdbUrls = append(tmdbUrls, fmt.Sprintf("%s%d", baseURL, searchMovie.Results[0].ID))

		// get poster
		if strings.HasSuffix(searchMovie.Results[0].PosterPath, "jpg") {
			posters = append(posters, tmdb.GetImageURL(searchMovie.Results[0].PosterPath, tmdb.W92))
		} else {
			posters = append(posters, "/img/no_movie_poster_found.jpg")
		}
	}
	// fmt.Println("POSTERS:", posters)
	// fmt.Println("URLS:", tmdbUrls)
	return posters, tmdbUrls
}

func unmarshallMovieTitles(data string) ([]string, []string) {
	var movies Movies
	err := json.Unmarshal([]byte(data), &movies)
	if err != nil {
		fmt.Println(err)
		return nil, nil // is this the best way to handle errors here?
	}

	var movieTitles []string
	for _, movie := range movies.Movies {
		movieTitles = append(movieTitles, fmt.Sprintf("%s (%d)", movie.Title, movie.Year))
	}

	var movieReasons []string
	for _, movie := range movies.Movies {
		movieReasons = append(movieReasons, movie.Reason)
	}
	
	// fmt.Println("TITLES:", movieTitles)
	// fmt.Println("REASONS:", movieReasons)
	return movieTitles, movieReasons
}
