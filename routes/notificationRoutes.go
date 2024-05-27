package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nosadvova/sportzone-backend/controllers"
)

func NotificationRoutes(router *gin.Engine) {
	router.GET("/notifications", controller.FetchNotifications())
	router.POST("/notifications", controller.CreateNotification())
	router.DELETE("/notifications/:notification_id", controller.DeleteNotification())
}
