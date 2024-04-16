package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"id"`
	First_name    string            `json:"first_name" validate:"required, min=2, max=100"`
	Last_name     string            `json:"last_name" validate:"required, min=2, max=100"`
	Password      string            `json:"password" validate:"required, min=6"` // Omit password from JSON response
	Email         string            `json:"email" validate:"email, required"`
	Phone         string            `json:"phone" validate:"required"`
	Created_at    time.Time         `json:"created_at"`
	User_id       string            `json:"user_id"`
}
