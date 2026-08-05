package controllers

import (
	"time"
	"weather-backend/config"
	"weather-backend/helpers" // Panggil helper yang sudah dibuat
	"weather-backend/middlewares"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := config.DB.Preload("Roles").Where("username = ?", input.Username).First(&user)
	if result.Error != nil || user.Password != input.Password {
		helpers.ResponseUnauthorized(c, gin.H{"error": "Username atau password salah!"})
		return
	}

	roleName := "user"
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     roleName,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middlewares.JwtKey)
	if err != nil {
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Respons sukses menggunakan helper
	helpers.ResponseOK(c, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
		"role":    roleName,
	})
}