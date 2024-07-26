package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/gofiber/fiber/v2"
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

type LogQuery struct {
	Query     string `json:"query"`
	Ip        string `json:"ip"`
	Time      string `json:"time"`
	UserAgent string `json:"user_agent"`
}

func Search(s3Client *s3.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")

		if query == "" {
			return c.Render("404", fiber.Map{
				"Message": "No query provided",
			}, "layouts/main")
		}

		go logQuery(c, s3Client)

		return c.Render("search", fiber.Map{
			"Title":       "FilmSearch - Search",
			"Canonical":   "https://filmsearch.xyz/search",
			"Robots":      "noindex, nofollow",
			"Description": "Results for: " + query,
			"Keywords":    "filmsearch, search, film, movie, discover, ai",
			"Query":       query,

		}, "layouts/main")
	}
}

func SearchResults(openaiClient *openai.Client, tmdbClient *tmdb.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")
		moviesJson, err := openaiMovieCompletion(openaiClient, query)
		if err != nil || moviesJson == "" {
			log.Println(err)
			return c.Render("404", fiber.Map{
				"Message": "No results for query: " + query,
			}, "layouts/main")
		}

		movieTitles, movieReasons := unmarshallMovieTitles(moviesJson)
		posters, tmdbUrls := getTmdbInfo(tmdbClient, movieTitles)

		// make sure all slices have the same length
		if !(len(movieTitles) == len(posters) && len(posters) == len(movieReasons) && len(movieReasons) == len(tmdbUrls)) {
			log.Println("Data length mismatch")
			return c.Render("404", fiber.Map{
				"Message": "Data length mismatch",
			}, "layouts/main")
		}

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
		log.Println(err)
		return nil, nil
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

// log query to aws s3
func logQuery(c *fiber.Ctx, s3Client *s3.Client) {
	query := LogQuery{
		Query:     c.Query("q"),
		Ip:        c.IP(),
		Time:      time.Now().Format(time.RFC3339),
		UserAgent: string(c.Request().Header.Peek("User-Agent")),
	}
	csvLine := fmt.Sprintf("%s,%s,%s,%s\n", query.Query, query.Ip, query.Time, query.UserAgent)

	getObjectOutput, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("filmsearch-query-log-654tizo86wufmowm8o34btuna4gukusw1b-s3alias"),
		Key:    aws.String("queries.csv"),
	})
	if err != nil {
		log.Printf("Error getting existing CSV from S3: %v", err)
	}

	var existingContent string
	if getObjectOutput.Body != nil {
		content, err := io.ReadAll(getObjectOutput.Body)
		if err != nil {
			log.Printf("Error reading existing CSV content: %v", err)
		} else {
			existingContent = string(content)
		}
		getObjectOutput.Body.Close()
	}

	newContent := existingContent + csvLine

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("filmsearch-query-log-654tizo86wufmowm8o34btuna4gukusw1b-s3alias"),
		Key:         aws.String("queries.csv"),
		Body:        strings.NewReader(newContent),
		ContentType: aws.String("text/csv"),
	})
	if err != nil {
		log.Printf("Error appending query to S3: %v", err)
	} else {
		log.Println("Logged query to S3")
	}
}