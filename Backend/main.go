package main

import (
	// "log"
	"project/configs"
	"project/routes"

	// dotenv "github.com/dsh2dsh/expx-dotenv"
	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.ConnectDatabase()
	app := fiber.New()

	userGroup := app.Group("/api")
	routes.UserRoute(userGroup)
	routes.TaskRoute(userGroup)

	app.Listen(":8080")

}
