package db

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"

	"go_meal_record/app/model"
)

var DB *gorm.DB
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
	DB, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	DB.AutoMigrate(&model.User{})
	// DbEngine.AutoMigrate(&model.Book{})
}
