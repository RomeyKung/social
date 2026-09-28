package main

import (
	"log"

	"github.com/RomeyKung/social/database"
	_ "github.com/RomeyKung/social/docs"
	"github.com/RomeyKung/social/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Social API
// @version 1.0
// @description This is a sample server for a social media application.
// @host localhost:5000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token

func main() {
	//load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .end file")
	}

	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	routes.SetupRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Listen(":5000")

}
