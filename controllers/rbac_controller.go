package controllers

import (
	"weather-backend/config"
	"weather-backend/helpers"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetRBACData mengambil seluruh data Role, Permission, dan relasinya dengan User
func GetRBACData(c *gin.Context) {
	var roles []models.Role
	var permissions []models.Permission
	var users []models.User

	// 1. Ambil semua Role beserta relasi Permissions-nya
	if err := config.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		config.Log.Error("Gagal mengambil data roles dari database", zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal mengambil data roles"})
		return
	}

	// 2. Ambil semua Permission yang tersedia
	if err := config.DB.Find(&permissions).Error; err != nil {
		config.Log.Error("Gagal mengambil data permissions dari database", zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal mengambil data permissions"})
		return
	}

	// 3. Ambil semua User beserta relasi Roles-nya (Sembunyikan password)
	if err := config.DB.Omit("password").Preload("Roles").Find(&users).Error; err != nil {
		config.Log.Error("Gagal mengambil data users dari database", zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal mengambil data users"})
		return
	}

	// Log sukses mengambil data RBAC
	config.Log.Info("Data RBAC berhasil diambil", 
		zap.Int("total_roles", len(roles)),
		zap.Int("total_permissions", len(permissions)),
		zap.Int("total_users", len(users)),
	)

	// 4. Return data menggunakan helper
	helpers.ResponseOK(c, gin.H{
		"message": "Data RBAC berhasil diambil",
		"data": gin.H{
			"roles":       roles,
			"permissions": permissions,
			"users":       users,
		},
	})
}