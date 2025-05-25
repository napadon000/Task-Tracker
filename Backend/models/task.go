package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title" validate:"required"`
	Description string        `json:"description" bson:"description"`
	Status      string        `json:"status" bson:"status" validate:"required,oneof=todo in_progress done"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at" `
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at" `
}

// type TaskStatus string

// const (
// 	StatusTodo       TaskStatus = "todo"
// 	StatusInProgress TaskStatus = "in_progress"
// 	StatusDone       TaskStatus = "done"
// )

var TaskValidate *validator.Validate

func init() {
	TaskValidate = validator.New()
}
