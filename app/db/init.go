package db

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"

	"go-meal-record/app/model"
)

var DB *gorm.DB
var err error

func init() {
	//circleciを利用する=githubを利用するということは.envは使えないので、開発環境だけ.env使って、本番環境(ecs)はタスク定義に環境変数を書く
	if os.Getenv("GO_ENVIRONMENT") == "development" {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalln(err)
		}
	}
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
