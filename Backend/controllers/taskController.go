package controllers

import (
	"context"
	"fmt"
	"os"
	"project/configs"
	"project/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var cliam jwt.MapClaims

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	cliam = token.Claims.(jwt.MapClaims)

	return c.Next()
}

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//validation
	if err := models.TaskValidate.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	result, err := configs.TaskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating task")
	}

	task.ID, _ = bson.ObjectIDFromHex(result.InsertedID.(bson.ObjectID).Hex())

	// Update the user's task array
	userID := cliam["id"].(string)
	fmt.Println("User ID from JWT:", userID)
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$push": bson.M{
			"tasks": task.ID,
		},
	}

	_, err = configs.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating user's tasks")
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func GetAllTasks(c *fiber.Ctx) error {
	userID := cliam["id"].(string)
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	filter := bson.M{"_id": id}
	var user models.User
	err = configs.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
			"id":    filter,
		})
	}

	tasks := []models.Task{}
	if len(user.Tasks) > 0 {
		cursor, err := configs.TaskCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": user.Tasks}})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching tasks")
		}

		if err := cursor.All(context.TODO(), &tasks); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error decoding tasks")
		}
	}

	return c.JSON(tasks)
}

func ChangeTaskStatus(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Task ID is required")
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if task.Status == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Task status is required")
	}

	id, err := bson.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     task.Status,
			"updated_at": time.Now(),
		},
	}

	result, err := configs.TaskCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating task status")
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Task not found")
	}

	return c.JSON(fiber.Map{"message": "Task status updated successfully"})
}

func DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Task ID is required")
	}

	id, err := bson.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	filter := bson.M{"_id": id}
	result, err := configs.TaskCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting task")
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Task not found")
	}

	// Remove the task ID from the user's tasks array
	userID := cliam["id"].(string)
	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	update := bson.M{
		"$pull": bson.M{
			"tasks": id,
		},
	}

	_, err = configs.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": userObjectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating user's tasks")
	}

	return c.JSON(fiber.Map{"message": "Task deleted successfully"})
}

func UpdateTaskDescription(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Task ID is required")
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if task.Description == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Task description is required")
	}

	id, err := bson.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"description": task.Description,
			"updated_at":  time.Now(),
		},
	}

	result, err := configs.TaskCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error updating task description")
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Task not found")
	}

	return c.JSON(fiber.Map{"message": "Task description updated successfully"})
}
