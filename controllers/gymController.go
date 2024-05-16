package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nosadvova/sportzone-backend/database"
	"github.com/nosadvova/sportzone-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var gymCollection *mongo.Collection = database.OpenCollection(database.Client, "gyms")

func CreateGym() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var gym models.Gym

		if err := c.BindJSON(&gym); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		trainerID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
		gym.Trainer_ID = trainerID
		gym.ID = primitive.NewObjectID()
		gym.Gym_Id = gym.ID.Hex()
		gym.Created_At = time.Now()

		gym.Sportsmen = []string{}
		gym.Trainings = []string{}

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": trainerID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user information"})
			return
		}
		if user.Personal_Information.Gym != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gym already exists for this trainer"})
			return
		}

		// Insert the gym into the database
		result, err := gymCollection.InsertOne(ctx, gym)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gym"})
			return
		}

		update := bson.M{"$set": bson.M{"personal_information.gym": gym.Gym_Id}}
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": trainerID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's gym information"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
	}
}

func GetGym() gin.HandlerFunc {
	return func(c *gin.Context) {
		gymId := c.Param("gym_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var gym models.Gym
		err := gymCollection.FindOne(ctx, bson.M{"gym_id": gymId}).Decode(&gym)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gym not found"})
			return
		}
		c.JSON(http.StatusOK, gym)
	}
}

func GetGyms() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <= 0 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		startIndex, _ := strconv.Atoi(c.Query("startIndex"))
		if startIndex < 0 {
			startIndex = (page - 1) * recordPerPage
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "_id", Value: "null"},
				}},
				{Key: "total_count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "data", Value: bson.D{
					{Key: "$push", Value: "$$ROOT"}}}}}}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "gyms", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := gymCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while getting gyms"})
			return
		}

		var gymsResult []bson.M
		if err = result.All(ctx, &gymsResult); err != nil {
			log.Fatal(err)
			return
		}

		if len(gymsResult) > 0 {
			c.JSON(http.StatusOK, gymsResult[0])
		} else {
			c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		}
	}
}

func FollowGym() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		gymID := c.Param("gym_id")
		userID := c.GetString("user_id")

		gymObjectID, err := primitive.ObjectIDFromHex(gymID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gym ID"})
			return
		}

		// Find the gym
		var gym models.Gym
		err = gymCollection.FindOne(ctx, bson.M{"_id": gymObjectID}).Decode(&gym)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
			return
		}

		if gym.Sportsmen == nil {
			gym.Sportsmen = []string{}
		}

		// Check if the user is already following the gym
		for _, id := range gym.Sportsmen {
			if id == userID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Already following this gym"})
				return
			}
		}

		// Add the user's ID to the gym's sportsmen array
		updateGym := bson.M{"$push": bson.M{"sportsmen": userID}}
		_, err = gymCollection.UpdateOne(ctx, bson.M{"_id": gymObjectID}, updateGym)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow gym"})
			return
		}

		// Update the user's gym field
		updateUser := bson.M{"$set": bson.M{"personal_information.gym": gymID}}
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": userID}, updateUser)
		if err != nil {
			// Rollback the previous operation if this fails
			gymCollection.UpdateOne(ctx, bson.M{"_id": gymObjectID}, bson.M{"$pull": bson.M{"sportsmen": userID}})
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's gym information"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully followed the gym"})
	}
}

func GetSportsmenForGym() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		gymID := c.Param("gym_id")

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
}
