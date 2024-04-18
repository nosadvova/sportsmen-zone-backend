package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Gym struct {
// 	ID          primitive.ObjectID `bson:"id"`
// 	Name        *string            `json:"name" validation:"required,min=2,max=20"`
// 	Description *string            `json:"description" validation:"required,min=2,max=200"`
// 	Created_At  time.Time          `json:"created_at"`
// 	Trainer     *User              `json:"trainer"`
// 	Sportsmen   []User             `json:"sportsmen"`
// }

type Gym struct {
	ID          primitive.ObjectID   `bson:"id"`
	Name        string               `json:"name" validation:"required"`
	Description string               `json:"description"`
	Location    Location             `json:"location"`
	Trainer_ID  primitive.ObjectID   `json:"trainer_id" validation:"required"`
	Sportsmen   []primitive.ObjectID `json:"sportsmen"`
	Trainings   []primitive.ObjectID `json:"trainings"`
	Created_At  time.Time            `json:"created_at"`
}

type Location struct {
	City            string `json:"city"`
	District        string `json:"district"`
	Street          string `json:"street"`
	Building_Number string `json:"building_number"`
}
