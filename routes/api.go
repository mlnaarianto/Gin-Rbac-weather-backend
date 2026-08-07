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

		// ---> RUTE RBAC (read-only, boleh tetap ikut manage-sensors atau nanti dipisah juga) <---
		adminRoutes.GET("/rbac", controllers.GetRBACData)
	}

	// ---> RUTE UPDATE ROLE USER - permission "manage-user-roles" <---
	userRoleRoutes := r.Group("/api/admin/users")
	userRoleRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("manage-user-roles"))
	{
		userRoleRoutes.PUT("/:id/roles", controllers.UpdateUserRoles)
	}

	// ---> RUTE MANAJEMEN ROLE (CRUD role) - permission "manage-roles" <---
	roleRoutes := r.Group("/api/admin/roles")
	roleRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("manage-roles"))
	{
		roleRoutes.POST("", controllers.CreateRole)
		roleRoutes.PUT("/:id", controllers.UpdateRole)
		roleRoutes.DELETE("/:id", controllers.DeleteRole)
	}

	// ---> RUTE MANAJEMEN PERMISSION ROLE (assign permission ke role) - permission "manage-role-permissions" <---
	rolePermRoutes := r.Group("/api/admin/roles")
	rolePermRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("manage-role-permissions"))
	{
		rolePermRoutes.PUT("/:id/permissions", controllers.UpdateRolePermissions)
	}
}