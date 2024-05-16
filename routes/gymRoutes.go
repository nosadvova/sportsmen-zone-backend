package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nosadvova/sportzone-backend/controllers"
)

func GymRoutes(router *gin.Engine) {
	router.POST("/gym", controller.CreateGym())
	router.GET("/gym", controller.GetGyms())
	router.GET("/gym/:gym_id", controller.GetGym())
	router.GET("/gym/sportsmen", controller.GetSportsmenForGym())
	router.POST("/gym/:gym_id/follow", controller.FollowGym())
}
