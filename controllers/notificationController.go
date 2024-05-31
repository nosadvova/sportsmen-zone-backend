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
		if trainerID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}

		trainerObjectID, err := primitive.ObjectIDFromHex(trainerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Trainer ID"})
			return
		}

		log.Print(notification.Gym_ID)

		gymObjectID := notification.Gym_ID

		var gym models.Gym
		err = gymCollection.FindOne(ctx, bson.M{"_id": gymObjectID, "trainer_id": trainerObjectID}).Decode(&gym)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Trainer does not own this gym"})
			return
		}

		notification.ID = primitive.NewObjectID()
		notification.Created_At = time.Now()
		notification.Trainer_ID = trainerObjectID

		_, err = notificationCollection.InsertOne(ctx, notification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
			return
		}

		_, err = gymCollection.UpdateOne(ctx, bson.M{"_id": gymObjectID}, bson.M{"$push": bson.M{"notifications": notification.ID}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gym with notification"})
			return
		}

		cursor, err := userCollection.Find(ctx, bson.M{"personal_information.gym": notification.Gym_ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer cursor.Close(ctx)

		var users []models.User
		if err = cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
			return
		}

		for _, user := range users {
			_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$push": bson.M{"notifications": notification.ID}})
			if err != nil {
				log.Printf("Failed to update user %s with notification: %v", user.ID.Hex(), err)
			}
		}

		c.JSON(http.StatusCreated, notification)
	}
}

func FetchNotifications() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))

		// Find the user
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Check if the user is associated with a gym
		if user.Personal_Information.Gym == nil {
			c.JSON(http.StatusOK, gin.H{"notifications": []models.Notification{}})
			return
		}

		// Convert Gym ID to ObjectID
		gymID, err := primitive.ObjectIDFromHex(*user.Personal_Information.Gym)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Gym ID"})
			return
		}

		var gym models.Gym
		err = gymCollection.FindOne(ctx, bson.M{"_id": gymID}).Decode(&gym)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gym not found"})
			return
		}

		// Fetch notifications using the notification IDs from the gym
		cursor, err := notificationCollection.Find(ctx, bson.M{"_id": bson.M{"$in": gym.Notifications}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
			return
		}
		defer cursor.Close(ctx)

		var notifications []models.Notification
		if err = cursor.All(ctx, &notifications); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode notifications"})
			return
		}

		c.JSON(http.StatusOK, notifications)
	}
}

func DeleteNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notificationID := c.Param("notification_id")
		trainerID := c.GetString("user_id")
		if trainerID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		notificationObjectID, err := primitive.ObjectIDFromHex(notificationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Notification ID"})
			return
		}

		trainerObjectID, err := primitive.ObjectIDFromHex(trainerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Trainer ID"})
			return
		}

		log.Printf("Deleting Notification ID: %s by Trainer ID: %s", notificationObjectID.Hex(), trainerObjectID.Hex())

		// Verify that the trainer created the notification
		var notification models.Notification
		err = notificationCollection.FindOne(ctx, bson.M{"_id": notificationObjectID, "trainer_id": trainerObjectID}).Decode(&notification)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Trainer did not create this notification"})
			return
		}

		// Remove the notification from the notifications collection
		_, err = notificationCollection.DeleteOne(ctx, bson.M{"_id": notificationObjectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
			return
		}

		// Remove the notification ID from the corresponding gym document
		_, err = gymCollection.UpdateOne(ctx, bson.M{"_id": notification.Gym_ID}, bson.M{"$pull": bson.M{"notifications": notificationObjectID}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove notification from gym"})
			return
		}

		// Optionally: Remove the notification ID from the user documents
		_, err = userCollection.UpdateMany(ctx, bson.M{}, bson.M{"$pull": bson.M{"notifications": notificationObjectID}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove notification from users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
	}
}
