package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nosadvova/sportzone-backend/controllers"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("users/login", controller.Login())
	router.POST("users/register", controller.Register())
	router.POST("users/refresh-token", controller.RefreshToken())
}
