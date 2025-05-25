package controllers

import (
	"context"
	"os"
	"project/configs"
	"project/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// func GetUserByID(c *fiber.Ctx) error {
// 	configs.UserCollection := configs.MongoClient.Database("tasktracker").Collection("user")

// 	userID, err := bson.ObjectIDFromHex(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
// 	}

// 	filter := bson.M{"_id": userID}

// 	var user models.User
// 	err = configs.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
// 	if err != nil {

// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": err.Error(),
// 			"id":    filter,
// 		})

// 	}
// 	return c.JSON(user)
// }

func CreateUser(c *fiber.Ctx) error {
	// configs.UserCollection := configs.MongoClient.Database("tasktracker").Collection("user")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//validation
	// models.UserValidate.RegisterValidation("username_format", models.ValidateUsername)

	if err := models.UserValidate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}

	user.Password = string(hashedPassword)

	if user.Tasks == nil {
		user.Tasks = []bson.ObjectID{}
	}

	result, err := configs.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).SendString("Email or Username already exists")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user.ID, _ = bson.ObjectIDFromHex(result.InsertedID.(bson.ObjectID).Hex())

	return c.Status(fiber.StatusCreated).JSON(user)
}

func LoginUser(c *fiber.Ctx) error {
	// configs.UserCollection := configs.MongoClient.Database("tasktracker").Collection("user")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// validation
	// models.UserValidate.RegisterValidation("username_format", models.Void)
	if err := models.UserValidate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	filter := bson.M{"email": user.Email}
	var foundUser models.User
	err := configs.UserCollection.FindOne(context.TODO(), filter).Decode(&foundUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Password does not match")
	}

	//login successful, generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  foundUser.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 72).Unix(), // token valid for 72 hours
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(72 * time.Hour), // token valid for 72 hours
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    foundUser,
		"token":   t,
	})
}
