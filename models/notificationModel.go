package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id"`
	GymID     string             `json:"gym_id"`
	Type      string             `json:"type"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	CreatedAt time.Time          `json:"created_at"`
}
