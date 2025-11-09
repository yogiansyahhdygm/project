package config

import (
	"fmt"
	"os"
	"toko/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	// Railway biasanya butuh SSLMode=require agar koneksi aman
	// tapi kamu bisa ubah ke disable kalau error SSL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	fmt.Println("Connected to Railway PostgreSQL!")
	return db, nil

}

// Auto migrate all models
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Profil{}, &models.Kategori{}, &models.Produk{}, &models.Pesanan{})
}
