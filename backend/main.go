package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/purnasavitri/personal-expense-tracker/backend/config"
	"github.com/purnasavitri/personal-expense-tracker/backend/models"
	"github.com/purnasavitri/personal-expense-tracker/backend/routes"
)

func main() {
	// Inisialisasi Gin
	router := gin.Default()

	// 2. TERAPKAN MIDDLEWARE CORS
	// Ini mengizinkan frontend di localhost:5173 untuk berkomunikasi dengan backend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Hubungkan ke Database
	config.ConnectToDB()

	// Migrasi Database
	config.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{})

	// Daftarkan Rute API
	routes.SetupAPIRoutes(router)

	// Jalankan Server
	router.Run() // Secara default berjalan di :8080
}