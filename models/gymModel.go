package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gym struct {
	ID          primitive.ObjectID `bson:"id"`
	Name        *string            `json:"name" validation:"required,min=2,max=20"`
	Description *string            `json:"description" validation:"required,min=2,max=200"`
	Created_At  time.Time          `json:"created_at"`
	Trainer     *User              `json:"trainer"`
	Sportsmen   []User             `json:"sportsmen"`
}
