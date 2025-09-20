package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purnasavitri/personal-expense-tracker/backend/config"
	"github.com/purnasavitri/personal-expense-tracker/backend/models"
)

// --- Membuat Transaksi Baru ---
func CreateTransaction(c *gin.Context) {
	var body struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		CategoryID  uint    `json:"category_id"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Ambil user yang sedang login dari context (akan di-set oleh middleware nanti)
	user, _ := c.Get("user")

	// Buat objek transaksi
	transaction := models.Transaction{
		Description: body.Description,
		Amount:      body.Amount,
		Type:        body.Type,
		CategoryID:  body.CategoryID,
		UserID:      user.(models.User).ID, // Set UserID dari user yang login
	}

	// Simpan ke database
	result := config.DB.Create(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// --- Mengambil Semua Transaksi ---
func GetTransactions(c *gin.Context) {
	user, _ := c.Get("user")

	var transactions []models.Transaction
	// Ambil semua transaksi yang UserID-nya cocok dengan user yang sedang login
	config.DB.Where("user_id = ?", user.(models.User).ID).Find(&transactions)

	c.JSON(http.StatusOK, transactions)
}

// --- Mengambil Satu Transaksi Berdasarkan ID ---
func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("user")

	var transaction models.Transaction
	// Cari transaksi berdasarkan ID DAN pastikan pemiliknya adalah user yang sedang login
	result := config.DB.First(&transaction, "id = ? AND user_id = ?", id, user.(models.User).ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// --- Mengubah Transaksi ---
func UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("user")

	var transaction models.Transaction
	result := config.DB.First(&transaction, "id = ? AND user_id = ?", id, user.(models.User).ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var body struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		CategoryID  uint    `json:"category_id"`
	}
	c.Bind(&body)

	// Update field di objek transaksi
	config.DB.Model(&transaction).Updates(models.Transaction{
		Description: body.Description,
		Amount:      body.Amount,
		Type:        body.Type,
		CategoryID:  body.CategoryID,
	})

	c.JSON(http.StatusOK, transaction)
}

// --- Menghapus Transaksi ---
func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("user")

	// Cari dan hapus transaksi jika ID dan UserID-nya cocok
	result := config.DB.Delete(&models.Transaction{}, "id = ? AND user_id = ?", id, user.(models.User).ID)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.Status(http.StatusNoContent)
}