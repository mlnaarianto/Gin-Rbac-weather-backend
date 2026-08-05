package database

import "log"

// RunAllSeeders menjalankan semua seeder secara berurutan
func RunAllSeeders() {
	SeedRoles()
	SeedPermissions()
	SeedUsers()
	log.Println("Seluruh database seeder berhasil dijalankan!")
}