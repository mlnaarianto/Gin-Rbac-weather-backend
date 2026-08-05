package database

import (
	"fmt"
	"weather-backend/config"
	"weather-backend/models"
)

func SeedPermissions() {
	manageSensorsPerm := models.Permission{Name: "manage-sensors"}
	viewWeatherPerm := models.Permission{Name: "view-weather"}

	config.DB.FirstOrCreate(&manageSensorsPerm, models.Permission{Name: "manage-sensors"})
	config.DB.FirstOrCreate(&viewWeatherPerm, models.Permission{Name: "view-weather"})

	// Ambil role admin & user untuk dihubungkan
	var adminRole, userRole models.Role
	config.DB.Where("name = ?", "admin").First(&adminRole)
	config.DB.Where("name = ?", "user").First(&userRole)

	// Hubungkan relasi role_permissions
	config.DB.Model(&adminRole).Association("Permissions").Append(&manageSensorsPerm, &viewWeatherPerm)
	config.DB.Model(&userRole).Association("Permissions").Append(&viewWeatherPerm)

	fmt.Println("Seeder Permissions & Role_Permissions berhasil dijalankan!")
}