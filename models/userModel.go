package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// add personal information instead of all fields below ID
type User struct {
	ID                   primitive.ObjectID   `bson:"_id"`
	Personal_Information *PersonalInformation `json:"personal_information"`
	// User_Image    *string             `json:"user_image"`
	// First_Name    *string             `json:"first_name" validation:"required,min=2,max=20"`
	// Last_Name     *string             `json:"last_name" validation:"required,min=2,max=20"`
	// Password      *string             `json:"password" validation:"required,min=6"`
	// Email         *string             `json:"email" validate:"required,email"`
	// User_Type     *string             `json:"user_type" validate:"required,eq=Sportsman|eq=Trainer"`
	// Gym           *primitive.ObjectID `json:"gym"`
	// Created_At    time.Time           `json:"created_at"`
	// Token         *string             `json:"token"`
	// Refresh_Token *string             `json:"refresh_token"`
	// User_id       string              `json:"user_id"`
}
