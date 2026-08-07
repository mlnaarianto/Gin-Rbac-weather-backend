package main

import (
	"log"
	"weather-backend/config"
	"weather-backend/database"
	"weather-backend/middlewares"
	"weather-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// 0. Muat file .env di awal sebelum konfigurasi database atau seeder
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Gagal memuat file .env, menggunakan variabel environment sistem.")
	}

	// 1. Inisialisasi Zap Logger di config
	config.InitLogger()
	defer config.Log.Sync() // Memastikan sisa buffer log ter-flush saat aplikasi mati

	config.Log.Info("Aplikasi sedang diinisialisasi...")

	// 2. Koneksi Database
	config.ConnectDB()

	// 3. Jalankan Seluruh Seeder Secara Terpusat
	database.RunAllSeeders()

	r := gin.Default()

	// 4. Gunakan Middleware CORS
	r.Use(middlewares.CORSMiddleware())

	// 5. Daftarkan Rute API
	routes.SetupRouter(r)

	// 6. Jalankan Server
	config.Log.Info("Server berhasil berjalan", zap.String("port", ":8080"))
	r.Run(":8080")
}