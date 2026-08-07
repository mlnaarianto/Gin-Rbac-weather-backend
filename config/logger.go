package config

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	// Inisialisasi Zap Logger
	Log, err = zap.NewDevelopment() // Bisa diganti zap.NewProduction() jika sudah siap production
	if err != nil {
		panic("Gagal menginisialisasi Zap Logger: " + err.Error())
	}
}