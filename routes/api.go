package routes

import (
	"net/http"
	"weather-backend/controllers"
	"weather-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Rute Publik
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong! Weather API Batam is running smoothly."})
	})
	r.POST("/login", controllers.Login)

	// Rute Terproteksi dengan Middleware Autentikasi & Permission RBAC
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"message": "Selamat datang, " + username.(string) + "!",
			})
		})
	}

	// Rute Khusus Admin dengan Permission "manage-sensors"
	adminRoutes := r.Group("/api/admin")
	adminRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("manage-sensors"))
	{
		adminRoutes.GET("/dashboard-stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":              "Selamat datang di Panel Kontrol Admin BMKG Batam!",
				"status_server_ai":     "Aktif & Normal",
				"total_stasiun_sensor": 5,
			})
		})
	}
}