package controllers

import (
	"time"
	"weather-backend/config"
	"weather-backend/helpers"
	"weather-backend/middlewares"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"` // Bisa diisi username atau email
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Mendukung pencarian berdasarkan username ATAU email
	result := config.DB.Preload("Roles").Where("username = ? OR email = ?", input.Username, input.Username).First(&user)
	
	// Cek apakah user ditemukan di database
	if result.Error != nil {
		helpers.ResponseUnauthorized(c, gin.H{"error": "Username atau password salah!"})
		return
	}

	// Cek kecocokan password dengan hash bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
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