package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gym struct {
	ID          primitive.ObjectID   `bson:"_id"`
	Image       string               `json:"image"`
	Name        string               `json:"name" validation:"required"`
	Description string               `json:"description"`
	Location    Location             `json:"location"`
	Type        []string             `json:"type"`
	Trainer_ID  primitive.ObjectID   `json:"trainer_id" validation:"required"`
	Sportsmen   []primitive.ObjectID `json:"sportsmen"`
	Trainings   []primitive.ObjectID `json:"trainings"`
	Created_At  time.Time            `json:"created_at"`
	Gym_Id      string               `json:"gym_id"`
}

type Location struct {
	City            string `json:"city"`
	District        string `json:"district"`
	Street          string `json:"street"`
	Building_Number string `json:"building_number"`
}
