package model

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"go_meal_record/app/utils/common"
)

type Meal struct {
	gorm.Model
	ID         int    `gorm:"primary_key" json:"ID"`
	MealType   string `validate:"meal_type,required" json:"meal_type"` //custom validate
	User       User   `json:"user" validate:"required"`
	UserID     int
	Menus      []Menu         `json:"menus"`
	MealImages []MealImage    `json:"meal_images"`
	Date       datatypes.Date `json:"date" validate:"required"`
}

type Menu struct {
	gorm.Model
	ID     int    `gorm:"primary_key" json:"ID"`
	Name   string `validate:"required,max=255" json:"name"`
	Meal   Meal   `json:"menu" validate:"required"`
	MealID int
}

type MealImage struct {
	gorm.Model
	ID     int    `gorm:"primary_key" json:"ID"`
	File   string `validate:"required,max=255" json:"file"`
	Meal   Meal   `json:"menu" validate:"required"`
	MealID int
}

//filepathにIDを含めたいためaftercreateする
func (meal_image *MealImage) AfterCreate(db *gorm.DB) (err error) {
	id := strconv.Itoa(meal_image.ID)
	filename := meal_image.File
	meal_image.File = "app/static/meal/" + id + "/" + filename
	db.Save(&meal_image)
	return
}

//file削除する
func (meal_image *MealImage) BeforeDelete(db *gorm.DB) (err error) {
	id := strconv.Itoa(meal_image.ID)
	filepath := "app/static/meal/" + id
	if _, err := os.Stat(filepath); err == nil { //file or directory 存在確認
		// err = os.RemoveAll(filepath)
		err = errors.New("remove fail")
		if err != nil {
			log.Println("file(directory) remove failed")
			return err
		}
	}
	return
}

func CustomValidateMealType(fl validator.FieldLevel) bool {
	meal_type_choices := []string{"morning", "lunch", "dinner", "other"}
	if common.Contains(meal_type_choices, fl.Field().String()) {
		return true
	}
	return false
}

// func (m Meal) UnMarshalJSON() ([]byte, error) {
// 	type meal Meal // prevent recursion
// 	x := meal(m)
// 	date = x.Date

// 	return json.Marshal(x)
// }
