package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nosadvova/sportzone-backend/database"
	helper "github.com/nosadvova/sportzone-backend/helpers"
	"github.com/nosadvova/sportzone-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var foundUser models.User
		var authToken models.AuthToken

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
		}

		err := userCollection.FindOne(ctx, bson.M{"personal_information.email": user.Personal_Information.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		passwordIsValid, msg := helper.VerifyPassword(*user.Personal_Information.Password, *foundUser.Personal_Information.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Personal_Information.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Personal_Information.Email, *foundUser.Personal_Information.First_Name,
			*foundUser.Personal_Information.Last_Name, *foundUser.Personal_Information.User_Type, foundUser.Personal_Information.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.Personal_Information.User_id)
		err = userCollection.FindOne(ctx, bson.M{"personal_information.user_id": foundUser.Personal_Information.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		authToken.Token = token
		authToken.Refresh_Token = refreshToken
		authToken.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		authToken.Expires_At = authToken.Created_At.Add(time.Hour * 120)

		c.JSON(http.StatusOK, authToken)
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Print(err.Error())
			defer cancel()
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			log.Print(validationErr.Error())
			defer cancel()
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"personal_information.email": user.Personal_Information.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking user email"})
		}

		password := helper.HashPassword(*user.Personal_Information.Password)
		user.Personal_Information.Password = &password

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			log.Print(err.Error())
			return
		}

		user.Personal_Information.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.Personal_Information.User_id = user.ID.Hex()

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("error while inserting user: %s", insertErr.Error())
			log.Print(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, user)
		// c.Status(http.StatusOK)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"personal_information.user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <= 0 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		startIndex, _ := strconv.Atoi(c.Query("startIndex"))

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
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while getting users"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"personal_information.user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var requestToken models.AuthToken
		if err := c.BindJSON(&requestToken); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Print(err.Error())
			return
		}

		// Find user by refresh token, you might need an index on this field for efficiency
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"refresh_token": user.Personal_Information.Refresh_Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		if time.Now().After(requestToken.Expires_At) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
			return
		}

		// Generate new tokens
		newToken, newRefreshToken, _ := helper.GenerateAllTokens(*user.Personal_Information.Email, *user.Personal_Information.First_Name,
			*user.Personal_Information.Last_Name, *user.Personal_Information.User_Type, user.Personal_Information.User_id)
		helper.UpdateAllTokens(newToken, newRefreshToken, user.Personal_Information.User_id)

		update := bson.M{
			"$set": bson.M{
				"token":         newToken,
				"refresh_token": newRefreshToken,
				"token_expires": time.Now().Add(24 * time.Hour * 7), // Set token expiration, adjust as necessary
			},
		}
		_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": user.Personal_Information.User_id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user tokens"})
			return
		}

		// Return the new tokens
		c.JSON(http.StatusOK, gin.H{"token": newToken, "refresh_token": newRefreshToken})
	}
}
