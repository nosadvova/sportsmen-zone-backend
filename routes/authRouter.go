package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nosadvova/sportzone-backend/controllers"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("user/login", controller.Login())
	router.POST("user/register", controller.Register())
}
