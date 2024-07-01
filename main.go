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
	app.Static("/", "./static/public")

	app.Get("/", handlers.Index)
	app.Get("/search", handlers.Search(openaiClient, tmdbClient))
	app.Get("/about", handlers.About)
	app.Get("/api/openai", handlers.Openai(openaiClient))
	app.Use(handlers.NotFound)

	log.Println("Server starting on port", *port)
	log.Fatal(app.Listen(*port)) // default port: 8080
}
