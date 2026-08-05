package main

import (
	"weather-backend/config"
	"weather-backend/database"
	"weather-backend/middlewares"
	"weather-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Koneksi Database
	config.ConnectDB()

	// 2. Jalankan Seluruh Seeder Secara Terpusat
	database.RunAllSeeders()

	r := gin.Default()

	// 3. Gunakan Middleware CORS
	r.Use(middlewares.CORSMiddleware())

	// 4. Daftarkan Rute API
	routes.SetupRouter(r)

	// 5. Jalankan Server
	r.Run(":8080")
}