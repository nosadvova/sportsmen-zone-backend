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

var notificationCollection *mongo.Collection = database.OpenCollection(database.Client, "notifications")

func CreateNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var notification models.Notification
		if err := c.BindJSON(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		trainerID := c.GetString("user_id")
		gymObjectID, err := primitive.ObjectIDFromHex(notification.GymID.Hex())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Gym ID"})
			return
		}

		var gym models.Gym
		err = gymCollection.FindOne(ctx, bson.M{"_id": gymObjectID, "trainer_id": trainerID}).Decode(&gym)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Trainer does not own this gym"})
			log.Print(gymObjectID, trainerID)
			return
		}

		notification.ID = primitive.NewObjectID()
		notification.CreatedAt = time.Now()

		_, err = notificationCollection.InsertOne(ctx, notification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
			return
		}

		c.JSON(http.StatusCreated, notification)
	}
}

func FetchNotifications() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID := c.GetString("user_id")
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		if user.Personal_Information.Gym == nil {
			c.JSON(http.StatusNoContent, gin.H{"error": "Gym not found"})
			return
		}

		gymID, err := primitive.ObjectIDFromHex(*user.Personal_Information.Gym)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Gym ID"})
			return
		}

		cursor, err := notificationCollection.Find(ctx, bson.M{"gym_id": gymID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
			return
		}

		var notifications []models.Notification
		if err = cursor.All(ctx, &notifications); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode notifications"})
			return
		}

		c.JSON(http.StatusOK, notifications)
	}
}
