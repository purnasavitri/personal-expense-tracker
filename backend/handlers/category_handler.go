package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purnasavitri/personal-expense-tracker/backend/config"
	"github.com/purnasavitri/personal-expense-tracker/backend/models"
)

// --- Membuat Kategori Baru ---
func CreateCategory(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user, _ := c.Get("user")
	category := models.Category{Name: body.Name, UserID: user.(models.User).ID}

	result := config.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// --- Mengambil Semua Kategori Milik User ---
func GetCategories(c *gin.Context) {
	user, _ := c.Get("user")
	var categories []models.Category
	config.DB.Where("user_id = ?", user.(models.User).ID).Find(&categories)
	c.JSON(http.StatusOK, categories)
}