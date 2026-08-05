package database

import (
	"fmt"
	"weather-backend/config"
	"weather-backend/models"
)

func SeedPermissions() {
	// 1. Buat atau cari permissions terlebih dahulu
	manageSensorsPerm := models.Permission{Name: "manage-sensors"}
	viewWeatherPerm := models.Permission{Name: "view-weather"}

	config.DB.FirstOrCreate(&manageSensorsPerm, models.Permission{Name: "manage-sensors"})
	config.DB.FirstOrCreate(&viewWeatherPerm, models.Permission{Name: "view-weather"})

	// 2. Ambil role admin & user dengan penanganan model yang aman
	var adminRole models.Role
	if err := config.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		fmt.Println("Gagal menemukan role admin:", err)
		return
	}

	var userRole models.Role
	if err := config.DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
		fmt.Println("Gagal menemukan role user:", err)
		return
	}

	// 3. Hubungkan relasi role_permissions menggunakan Association GORM
	// Menggunakan Clear() terlebih dahulu opsional agar tidak terjadi duplikasi data saat seeder dijalankan ulang
	config.DB.Model(&adminRole).Association("Permissions").Replace(&manageSensorsPerm, &viewWeatherPerm)
	config.DB.Model(&userRole).Association("Permissions").Replace(&viewWeatherPerm)

	fmt.Println("Seeder Permissions & Role_Permissions berhasil dijalankan!")
}