package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/purnasavitri/personal-expense-tracker/backend/config"
	"github.com/purnasavitri/personal-expense-tracker/backend/models"
)

func RequireAuth(c *gin.Context) {
	// Ambil token dari Authorization header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Pisahkan "Bearer" dari token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader { // Jika tidak ada prefix "Bearer ", formatnya salah
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse claims"})
		return
	}

	// Cek apakah token sudah kedaluwarsa
	expFloat, ok := claims["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > expFloat {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}

	// Cari user di database berdasarkan ID dari token
	var user models.User
	config.DB.First(&user, claims["sub"])

	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Simpan data user ke dalam context request agar bisa digunakan oleh handler
	c.Set("user", user)

	// Lanjutkan ke handler berikutnya
	c.Next()
}