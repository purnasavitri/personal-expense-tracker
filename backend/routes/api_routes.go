package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/purnasavitri/personal-expense-tracker/backend/handlers"
	"github.com/purnasavitri/personal-expense-tracker/backend/middleware"
)

func SetupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Rute untuk autentikasi (tidak perlu dijaga)
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// Grup rute baru yang akan dilindungi oleh middleware RequireAuth
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth)
		{
			// Rute untuk Transaksi
			protected.POST("/transactions", handlers.CreateTransaction)
			protected.GET("/transactions", handlers.GetTransactions)
			protected.GET("/transactions/:id", handlers.GetTransactionByID)
			protected.PUT("/transactions/:id", handlers.UpdateTransaction)
			protected.DELETE("/transactions/:id", handlers.DeleteTransaction)
		}
	}
}