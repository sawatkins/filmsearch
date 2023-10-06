package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"

	"github.com/sawatkins/eureka-search/handlers"

)

var (
	port    = flag.String("port", ":8080", "Port to listen on")
	prefork = flag.Bool("prefork", false, "Enable prefork in Production")
	dev     = flag.Bool("dev", true, "Enable development mode")
)

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	// Connected with database
	//database.Connect()

	// Load AWS credentials
	// cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	// if err != nil {
	// 	log.Printf("Failed to load AWS configuration, %v", err)
	// }
	// s3Client := s3.NewFromConfig(cfg)
	// s3PresignClient := s3.NewPresignClient(s3Client)

	// Create a new engine
	engine := html.New("./views", ".html")
	if *dev {
		engine.Reload(true) 
		engine.Debug(true)  
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prefork, // go run app.go -prefork
		Views:   engine,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Setup static files
	app.Static("/", "./static/public")

	// Create a /api/v1 endpoint
	// v1 := app.Group("/api/v1")
	// userApis := v1.Group("/user")
	// userApis.Post("/createUser", handlers.CreateUser)

	// Routes
	app.Get("/", handlers.Index)
	app.Get("/faq", handlers.Faq)

	// Handle not founds
	app.Use(handlers.NotFound)

	log.Println("Server starting on port", *port)

	// Listen on port 8080
	log.Fatal(app.Listen(*port)) // go run app.go -port=:8080
}