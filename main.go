package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"

	"github.com/sawatkins/eureka-search/handlers"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file!")
	}
}

func main() {
	port := flag.String("port", ":8080", "Port to listen on")
	prefork := flag.Bool("prefork", false, "Enable prefork in Production")
	dev := flag.Bool("dev", true, "Enable development mode")
	flag.Parse()

	loadEnv()
	openaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	tmdbClient, _ := tmdb.Init(os.Getenv("TMDB_API_KEY"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	if err != nil {
		log.Printf("Failed to load AWS configuration, %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	engine := html.New("./views", ".html")
	if *dev {
		engine.Reload(true)
		engine.Debug(true)
	}

	app := fiber.New(fiber.Config{
		Prefork: *prefork,
		Views:   engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Static("/", "./static")

	app.Get("/", handlers.Index)
	app.Get("/search", handlers.Search(s3Client))
	app.Get("/search-results", handlers.SearchResults(openaiClient, tmdbClient))
	app.Get("/about", handlers.About)
	app.Use(handlers.NotFound)

	log.Println("Server starting on port", *port)
	log.Fatal(app.Listen(*port)) // default port: 8080
}
