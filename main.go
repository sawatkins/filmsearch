package main

import (
	"flag"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	tmdb "github.com/cyruzin/golang-tmdb"

	"github.com/sawatkins/eureka-search/handlers"
)

type Config struct {
	Port    string
	Prefork bool
	Dev     bool
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file!")
	}
}

func main() {
	config := Config{
		Port:    *flag.String("port", ":8080", "Port to listen on"),
		Prefork: *flag.Bool("prefork", false, "Enable prefork in Production"),
		Dev:     *flag.Bool("dev", true, "Enable development mode"),
	}
	flag.Parse()

	loadEnv()
	openaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	tmdbClient, _ := tmdb.Init(os.Getenv("TMDB_API_KEY"))

	engine := html.New("./views", ".html")
	if config.Dev {
		engine.Reload(true)
		engine.Debug(true)
	}

	app := fiber.New(fiber.Config{
		Prefork: config.Prefork,
		Views:   engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Static("/", "./static/public")

	app.Get("/", handlers.Index)
	app.Get("/search", handlers.Search(openaiClient, tmdbClient))
	app.Get("/about", handlers.About)
	app.Get("/api/openai", handlers.Openai(openaiClient))
	app.Use(handlers.NotFound)

	log.Println("Server starting on port", config.Port)
	log.Fatal(app.Listen(config.Port)) // default port: 8080
}
