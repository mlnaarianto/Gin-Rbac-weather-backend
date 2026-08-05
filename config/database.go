package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Ubah "Loc=Local" menjadi "loc=Local" (huruf 'l' kecil)
	dsn := "root:@tcp(127.0.0.1:3306)/weather_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}

	fmt.Println("Berhasil terhubung ke database MySQL (weather_db)!")
}