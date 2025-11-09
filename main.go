package main

import (
	"log"
	"os"

	"toko/config"
	"toko/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, relying on environment variables")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// run auto migrate
	config.Migrate(db)

	r := gin.Default()
	routes.Setup(r, db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
