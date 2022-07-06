package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *gorm.DB
var err error

func init() {
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable", postgres_user, postgres_password, postgres_db)
	db, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// DbEngine.AutoMigrate(&model.User{})
	// DbEngine.AutoMigrate(&model.Book{})
}
