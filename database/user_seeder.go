package database

import (
	"fmt"
	"weather-backend/config"
	"weather-backend/models"
)

func SeedUsers() {
	// 1. Cek atau buat user superadmin
	var adminUser models.User
	if err := config.DB.Where("username = ?", "superadmin").First(&adminUser).Error; err != nil {
		adminUser = models.User{
			Username: "superadmin",
			Password: "123",
		}
		if errCreate := config.DB.Create(&adminUser).Error; errCreate != nil {
			fmt.Println("Gagal membuat user superadmin:", errCreate)
			return
		}
	}

	// 2. Cek atau buat user warga_batam
	var regularUser models.User
	if err := config.DB.Where("username = ?", "warga_batam").First(&regularUser).Error; err != nil {
		regularUser = models.User{
			Username: "warga_batam",
			Password: "123",
		}
		if errCreate := config.DB.Create(&regularUser).Error; errCreate != nil {
			fmt.Println("Gagal membuat user warga_batam:", errCreate)
			return
		}
	}

	// 3. Ambil role admin & user dengan validasi error yang aman
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

	// 4. Hubungkan relasi user_roles (tabel pivot) menggunakan .Replace()
	if err := config.DB.Model(&adminUser).Association("Roles").Replace(&adminRole); err != nil {
		fmt.Println("Gagal menghubungkan role admin ke superadmin:", err)
		return
	}

	if err := config.DB.Model(&regularUser).Association("Roles").Replace(&userRole); err != nil {
		fmt.Println("Gagal menghubungkan role user ke warga_batam:", err)
		return
	}

	fmt.Println("Seeder Users & User_Roles berhasil dijalankan!")
}