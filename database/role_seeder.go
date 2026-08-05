package database

import (
	"fmt"
	"weather-backend/config"
	"weather-backend/models"
)

func SeedRoles() {
	adminRole := models.Role{Name: "admin"}
	userRole := models.Role{Name: "user"}

	config.DB.FirstOrCreate(&adminRole, models.Role{Name: "admin"})
	config.DB.FirstOrCreate(&userRole, models.Role{Name: "user"})

	fmt.Println("Seeder Roles berhasil dijalankan!")
}