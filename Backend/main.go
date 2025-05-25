package main

import (
	// "log"
	"project/configs"
	"project/routes"

	// dotenv "github.com/dsh2dsh/expx-dotenv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	configs.ConnectDatabase()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	userGroup := app.Group("/api")
	routes.UserRoute(userGroup)
	routes.TaskRoute(userGroup)

	app.Listen(":8080")

}
