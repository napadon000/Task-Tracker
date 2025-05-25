package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID       bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Username string          `json:"username" bson:"username" validate:"required,min=3,max=20,username_format"`
	Email    string          `json:"email" bson:"email" validate:"required,email,max=100"`
	Password string          `json:"password" bson:"password" validate:"required,min=8,max=128,password_format"`
	Tasks    []bson.ObjectID `json:"tasks" bson:"tasks" `
}

// // NewUser creates a new user with default values
// func NewUser() User {
// 	return User{
// 		Tasks: []bson.ObjectID{},
// 	}
// }

// Custom username validation
func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Username can contain letters, numbers, underscores, and hyphens
	// Must start with a letter
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_-]*$`, username)
	return matched
}

// Complex password validation
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// Check for required character types
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func Void(fl validator.FieldLevel) bool {
	// This function is intentionally left empty to avoid unused variable error
	return true
}

var UserValidate *validator.Validate

// Register custom validations
func init() {
	UserValidate = validator.New()
	UserValidate.RegisterValidation("password_format", ValidatePassword)
}
