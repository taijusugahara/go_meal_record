package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"

	"go-meal-record/app/model"
)

var DB *gorm.DB
var err error

func SettingDb() {
	if os.Getenv("GO_ENVIRONMENT") == "test" {
		SettingTestDb()
	} else {
		postgres_host := "postgres"
		if os.Getenv("GO_ENVIRONMENT") == "production" {
			postgres_host = os.Getenv("POSTGRES_ENDPOINT")
		}
		postgres_user := os.Getenv("POSTGRES_USER")
		postgres_password := os.Getenv("POSTGRES_PASSWORD")
		postgres_db := os.Getenv("POSTGRES_DB")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", postgres_host, postgres_user, postgres_password, postgres_db)
		DB, err = gorm.Open(
			postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}

		DB.AutoMigrate(&model.User{})
		DB.AutoMigrate(&model.Meal{})
		DB.AutoMigrate(&model.Menu{})
		DB.AutoMigrate(&model.MealImage{})
	}
}
