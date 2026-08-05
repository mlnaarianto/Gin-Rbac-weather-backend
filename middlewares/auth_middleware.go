package middlewares

import (
	"net/http"
	"os"
	"strings"
	"weather-backend/config"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Mengambil kunci rahasia JWT secara dinamis dari file .env
var JwtKey = []byte(getJwtSecret())

func getJwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "default-fallback-secret" // Cadangan jika .env tidak terbaca
	}
	return secret
}

// 1. Middleware untuk verifikasi token JWT (Autentikasi)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan, silakan login dulu"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
			c.Abort()
			return
		}

		// Simpan username ke context
		c.Set("username", claims["username"])
		c.Next()
	}
}

// 2. Middleware RBAC Berbasis Permission yang Lebih Aman
func RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		err := config.DB.Preload("Roles").Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User tidak ditemukan di database"})
			c.Abort()
			return
		}

		// Periksa apakah user memiliki permission yang dimaksud
		hasPermission := false
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Name == requiredPermission {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Anda tidak memiliki izin (permission) tersebut."})
			c.Abort()
			return
		}

		c.Next()
	}
}