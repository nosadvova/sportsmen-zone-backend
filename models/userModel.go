package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"id"`
	User_Image    *string            `json:"user_image"`
	First_Name    *string            `json:"first_name" validation:"required,min=2,max=20"`
	Last_Name     *string            `json:"last_name" validation:"required,min=2,max=20"`
	Password      *string            `json:"password" validation:"required,min=6"`
	Email         *string            `json:"email" validate:"required,email"`
	User_Type     *string            `json:"user_type" validate:"required,eq=Sportsman|eq=Trainer"`
	Created_At    time.Time          `json:"created_at"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	User_id       string             `json:"user_id"`
}
