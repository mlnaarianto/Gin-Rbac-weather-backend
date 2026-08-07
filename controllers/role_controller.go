package controllers

import (
	"weather-backend/config"
	"weather-backend/helpers"
	"weather-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpdateRolePermissions mengganti/menimpa seluruh permission milik sebuah role
// Endpoint: PUT /api/admin/roles/:id/permissions
// Body: { "permission_ids": ["uuid-1", "uuid-2", ...] }
func UpdateRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")

	var input struct {
		PermissionIDs []string `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Log.Warn("Gagal bind JSON saat update permission role",
			zap.String("role_id", roleIDStr),
			zap.Error(err),
		)
		helpers.ResponseBadRequest(c, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	// Parse string UUID role dari URL parameter
	parsedRoleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": "Format UUID role tidak valid"})
		return
	}

	// 1. Cari role berdasarkan UUID (konversi ke models.UUIDBinary)
	var role models.Role
	if err := config.DB.First(&role, "id = ?", models.UUIDBinary(parsedRoleID)).Error; err != nil {
		config.Log.Error("Role tidak ditemukan di database",
			zap.String("role_id", roleIDStr),
			zap.Error(err),
		)
		helpers.ResponseNotFound(c, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	// 2. Konversi string UUID permission ke models.UUIDBinary
	var permissionUUIDs []models.UUIDBinary
	for _, idStr := range input.PermissionIDs {
		if parsedPermID, err := uuid.Parse(idStr); err == nil {
			permissionUUIDs = append(permissionUUIDs, models.UUIDBinary(parsedPermID))
		}
	}

	// 3. Cari data Permission berdasarkan array UUIDBinary
	var permissions []models.Permission
	if len(permissionUUIDs) > 0 {
		config.DB.Find(&permissions, "id IN ?", permissionUUIDs)
	}

	// 4. Timpa/Ganti relasi many-to-many permissions milik role
	if err := config.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		config.Log.Error("Gagal eksekusi GORM association replace permission",
			zap.String("role_id", roleIDStr),
			zap.Error(err),
		)
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal memperbarui permission role"})
		return
	}

	config.Log.Info("Permission role berhasil diperbarui",
		zap.String("role_id", roleIDStr),
		zap.String("role_name", role.Name),
		zap.Any("new_permission_ids", input.PermissionIDs),
	)

	helpers.ResponseOK(c, gin.H{
		"message": "Permission role berhasil diperbarui",
	})
}

// CreateRole membuat role baru
// Endpoint: POST /api/admin/roles
// Body: { "name": "editor" }
func CreateRole(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	role := models.Role{
		Name: input.Name,
	}

	if err := config.DB.Create(&role).Error; err != nil {
		config.Log.Error("Gagal membuat role baru", zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal membuat role, kemungkinan nama sudah dipakai"})
		return
	}

	config.Log.Info("Role baru berhasil dibuat", zap.String("role_name", role.Name))

	helpers.ResponseOK(c, gin.H{
		"message": "Role berhasil dibuat",
		"data":    role,
	})
}

// UpdateRole mengubah nama role
// Endpoint: PUT /api/admin/roles/:id
// Body: { "name": "editor-baru" }
func UpdateRole(c *gin.Context) {
	roleIDStr := c.Param("id")

	parsedRoleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": "Format UUID role tidak valid"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, "id = ?", models.UUIDBinary(parsedRoleID)).Error; err != nil {
		helpers.ResponseNotFound(c, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	role.Name = input.Name
	if err := config.DB.Save(&role).Error; err != nil {
		config.Log.Error("Gagal memperbarui role", zap.String("role_id", roleIDStr), zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal memperbarui role"})
		return
	}

	helpers.ResponseOK(c, gin.H{
		"message": "Role berhasil diperbarui",
		"data":    role,
	})
}

// DeleteRole menghapus role
// Endpoint: DELETE /api/admin/roles/:id
func DeleteRole(c *gin.Context) {
	roleIDStr := c.Param("id")

	parsedRoleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		helpers.ResponseBadRequest(c, gin.H{"error": "Format UUID role tidak valid"})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, "id = ?", models.UUIDBinary(parsedRoleID)).Error; err != nil {
		helpers.ResponseNotFound(c, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	// Bersihkan relasi many-to-many sebelum hapus, hindari FK constraint error
	if err := config.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		config.Log.Error("Gagal membersihkan relasi permission sebelum hapus role", zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal menghapus role"})
		return
	}

	if err := config.DB.Delete(&role).Error; err != nil {
		config.Log.Error("Gagal menghapus role", zap.String("role_id", roleIDStr), zap.Error(err))
		helpers.ResponseInternalError(c, gin.H{"error": "Gagal menghapus role"})
		return
	}

	config.Log.Info("Role berhasil dihapus", zap.String("role_id", roleIDStr))

	helpers.ResponseOK(c, gin.H{
		"message": "Role berhasil dihapus",
	})
}