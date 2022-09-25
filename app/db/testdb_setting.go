package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"

	"go-meal-record/app/model"
)

func SettingTestDb() {
	postgres_host := "test_postgres"
	is_circleci_test := os.Getenv("IS_CIRCLECI_TEST")
	if is_circleci_test != "true" {
		postgres_host = "localhost" //circleciのimageの場合localhostになるらしい。
	}
	postgres_user := os.Getenv("TEST_POSTGRES_USER")
	postgres_password := os.Getenv("TEST_POSTGRES_PASSWORD")
	postgres_db := os.Getenv("TEST_POSTGRES_DB")
	//portが5433じゃなくて5432なのは外部からの接続ではなく内部からの接続になるらしい
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", postgres_host, postgres_user, postgres_password, postgres_db)
	DB, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	//dataリセットするため1度table落とす
	DB.Migrator().DropTable(&model.User{})
	DB.Migrator().DropTable(&model.Meal{})
	DB.Migrator().DropTable(&model.Menu{})
	DB.Migrator().DropTable(&model.MealImage{})
	//再度table作成
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Meal{})
	DB.AutoMigrate(&model.Menu{})
	DB.AutoMigrate(&model.MealImage{})
}
