package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID         primitive.ObjectID `bson:"_id"`
	Gym_ID     string             `json:"gym_id"`
	Trainer_ID primitive.ObjectID `json:"trainer_id"`
	Type       string             `json:"type"`
	Title      string             `json:"title"`
	Message    string             `json:"message"`
	Created_At time.Time          `json:"created_at"`
}
