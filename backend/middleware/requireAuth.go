package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/purnasavitri/personal-expense-tracker/backend/config"
	"github.com/purnasavitri/personal-expense-tracker/backend/models"
)

func RequireAuth(c *gin.Context) {
	// Ambil token dari Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
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
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Cek apakah token sudah kedaluwarsa
	expFloat, ok := claims["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > expFloat {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Cari user di database berdasarkan ID dari token
	var user models.User
	config.DB.First(&user, claims["sub"])

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Simpan data user ke dalam context request agar bisa digunakan oleh handler
	c.Set("user", user)

	// Lanjutkan ke handler berikutnya
	c.Next()
}