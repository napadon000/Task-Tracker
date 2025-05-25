package routes

import (
	"project/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app fiber.Router) {
	// app.Get("/user/:id", controllers.GetUserByID)
	app.Post("/user/register", controllers.CreateUser)
	app.Post("/user/login", controllers.LoginUser)
}
