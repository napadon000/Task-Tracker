package routes

import (
	"project/controllers"

	"github.com/gofiber/fiber/v2"
)

func TaskRoute(app fiber.Router) {
	app.Use(controllers.AuthRequired) // Apply authentication middleware to all task routes
	app.Post("/task/create", controllers.CreateTask)
	app.Get("/task/gettasks", controllers.GetAllTasks)
	app.Patch("/task/updatestatus/:id", controllers.ChangeTaskStatus)
	app.Patch("/task/updatetask/:id", controllers.UpdateTaskDescription)
	app.Delete("/task/deletetask/:id", controllers.DeleteTask)
}
