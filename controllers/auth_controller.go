package controllers

import (
	"time"
	"weather-backend/config"
	"weather-backend/helpers"
	"weather-backend/middlewares"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Log.Warn("Gagal bind JSON saat proses login", zap.Error(err))
		helpers.ResponseBadRequest(c, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// 1. Pastikan Preload("Roles") terpanggil dengan benar saat mencari user
	result := config.DB.Preload("Roles").Where("username = ? OR email = ?", input.Username, input.Username).First(&user)

	if result.Error != nil {
		config.Log.Warn("Login gagal: Username atau email tidak ditemukan", zap.String("input_username", input.Username))
		helpers.ResponseUnauthorized(c, gin.H{"error": "Username atau password salah!"})
		return
	}

	// 2. Cek kecocokan password dengan hash bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		config.Log.Warn("Login gagal: Password salah", zap.String("username", user.Username))
		helpers.ResponseUnauthorized(c, gin.H{"error": "Username atau password salah!"})
		return
	}

	// 3. Tentukan role secara cerdas (Cek apakah user punya role "admin")
	roleName := "user"
	for _, role := range user.Roles {
		if role.Name == "admin" {
			roleName = "admin"
			break
		}
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
		config.Log.Error("Gagal membuat token JWT saat login", zap.String("username", user.Username), zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Log sukses login
	config.Log.Info("Login berhasil",
		zap.String("username", user.Username),
		zap.String("role", roleName),
	)

	// Respons sukses
	helpers.ResponseOK(c, gin.H{
		"message":  "Login berhasil!",
		"token":    tokenString,
		"username": user.Username,
		"name":     user.Name,
		"role":     roleName,
	})
}