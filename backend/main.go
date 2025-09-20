// file: backend/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/purnasavitri/personal-expense-tracker/backend/config"
	"github.com/purnasavitri/personal-expense-tracker/backend/models"
	"github.com/purnasavitri/personal-expense-tracker/backend/routes"
)

func main() {
	// Inisialisasi Gin
	router := gin.Default()

	// Hubungkan ke Database
	config.ConnectToDB()

	// Migrasi Database (membuat tabel secara otomatis)
	config.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{})

	// Daftarkan Rute API
	routes.SetupAPIRoutes(router)

	// Jalankan Server
	router.Run() // Secara default berjalan di :8080
}