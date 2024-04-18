package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/nosadvova/sportzone-backend/controllers"
	"github.com/nosadvova/sportzone-backend/middleware"
)

func UserRoutes(router *gin.Engine) {
	router.Use(middleware.Authenticate())

	router.GET("/users", controller.GetUsers())
	router.GET("/users/:user_id", controller.GetUser())
}
