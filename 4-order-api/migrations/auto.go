package main

import (
	"dz4/internal/models"
	"dz4/internal/user"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // открываем соединение, чтобы подключиться в бд
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Product{}, &user.User{}, &models.Order{}); err != nil {
		log.Fatalln("failed to migrate:", err)
	} else {
		log.Println("Migration successful")
	}
}
