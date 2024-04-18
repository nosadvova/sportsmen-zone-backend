package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nosadvova/sportzone-backend/database"
	"github.com/nosadvova/sportzone-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var gymCollection *mongo.Collection = database.OpenCollection(database.Client, "gyms")

func CreateGym(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var gym models.Gym
	if err := c.BindJSON(&gym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trainerID, _ := primitive.ObjectIDFromHex(c.GetString("userID")) // Placeholder for actual user ID extraction
	gym.Trainer_ID = trainerID
	gym.Created_At = time.Now()

	// Insert the gym into the database
	result, err := gymCollection.InsertOne(ctx, gym)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gym"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func GetSportsmenForGym(c *gin.Context, gymID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gymObjectID, err := primitive.ObjectIDFromHex(gymID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gym ID"})
		return
	}

	matchStage := bson.D{{Key: "$match", Value: bson.M{"_id": gymObjectID}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "sportsmen",
		"localField":   "sportsmen",
		"foreignField": "_id",
		"as":           "sportsmenDetails",
	}}}

	cursor, err := gymCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sportsmen"})
		log.Print("Failed to retrieve sportsmen")
		return
	}

	var gymsWithSportsmen []bson.M
	if err = cursor.All(ctx, &gymsWithSportsmen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode sportsmen data"})
		log.Print("Failed to decode sportsmen data")
		return
	}

	c.JSON(http.StatusOK, gymsWithSportsmen)
}
