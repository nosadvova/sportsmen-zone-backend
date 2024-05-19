package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                   primitive.ObjectID   `bson:"_id"`
	Personal_Information *PersonalInformation `json:"personal_information"`
}
