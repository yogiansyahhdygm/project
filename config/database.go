package config

import (
	"fmt"
	"os"
	"toko/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSL_MODE")
	if ssl == "" {
		ssl = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta", host, user, pass, name, port, ssl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Auto migrate all models
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Profil{}, &models.Kategori{}, &models.Produk{}, &models.Pesanan{})
}
