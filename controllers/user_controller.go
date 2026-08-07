package controllers

import (
	"weather-backend/config"
	"weather-backend/helpers"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func UpdateUserRoles(c *gin.Context) {
	userIDStr := c.Param("id")

	var input struct {
		RoleIDs []string `json:"role_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Log.Warn("Gagal bind JSON saat update role",
			zap.String("user_id", userIDStr),
			zap.Error(err),
		)
		helpers.ResponseBadRequest(c, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	// Parse string UUID dari URL parameter ke tipe uuid.UUID
	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": "Format UUID user tidak valid"})
		return
	}

	// 1. Cari user berdasarkan UUID objek
	// PENTING: konversi ke models.UUIDBinary supaya cocok dengan tipe kolom binary(16) di DB
	var user models.User
	if err := config.DB.First(&user, "id = ?", models.UUIDBinary(parsedUserID)).Error; err != nil {
		config.Log.Error("Pengguna tidak ditemukan di database",
			zap.String("user_id", userIDStr),
			zap.Error(err),
		)
		helpers.ResponseNotFound(c, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// 2. Konversi string UUID role ke objek models.UUIDBinary (bukan uuid.UUID biasa)
	var roleUUIDs []models.UUIDBinary
	for _, idStr := range input.RoleIDs {
		if parsedRoleID, err := uuid.Parse(idStr); err == nil {
			roleUUIDs = append(roleUUIDs, models.UUIDBinary(parsedRoleID))
		}
	}

	// 3. Cari data Role berdasarkan array UUIDBinary
	var roles []models.Role
	if len(roleUUIDs) > 0 {
		config.DB.Find(&roles, "id IN ?", roleUUIDs)
	}

	// 4. Timpa/Ganti relasi many-to-many roles milik user
	if err := config.DB.Model(&user).Association("Roles").Replace(&roles); err != nil {
		config.Log.Error("Gagal eksekusi GORM association replace role",
			zap.String("user_id", userIDStr),
			zap.Error(err),
		)
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal memperbarui role pengguna"})
		return
	}

	config.Log.Info("Role pengguna berhasil diperbarui",
		zap.String("user_id", userIDStr),
		zap.String("username", user.Username),
		zap.Any("new_role_ids", input.RoleIDs),
	)

	helpers.ResponseOK(c, gin.H{
		"message": "Role pengguna berhasil diperbarui",
	})
}