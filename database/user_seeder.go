package database

import (
	"fmt"
	"weather-backend/config"
	"weather-backend/models"
)

func SeedUsers() {
	var adminUser models.User
	if err := config.DB.Where("username = ?", "superadmin").First(&adminUser).Error; err != nil {
		adminUser = models.User{
			Username: "superadmin",
			Password: "123",
		}
		config.DB.Create(&adminUser)
	}

	var regularUser models.User
	if err := config.DB.Where("username = ?", "warga_batam").First(&regularUser).Error; err != nil {
		regularUser = models.User{
			Username: "warga_batam",
			Password: "123",
		}
		config.DB.Create(&regularUser)
	}

	// Ambil role untuk dihubungkan ke user
	var adminRole, userRole models.Role
	config.DB.Where("name = ?", "admin").First(&adminRole)
	config.DB.Where("name = ?", "user").First(&userRole)

	// Hubungkan relasi user_roles (tabel pivot)
	config.DB.Model(&adminUser).Association("Roles").Append(&adminRole)
	config.DB.Model(&regularUser).Association("Roles").Append(&userRole)

	fmt.Println("Seeder Users & User_Roles berhasil dijalankan!")
}