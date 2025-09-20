package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/purnasavitri/personal-expense-tracker/backend/handlers"
)

func SetupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
	}
}