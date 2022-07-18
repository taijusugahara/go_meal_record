package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"strconv"
)

type Meal struct {
	gorm.Model
	ID         int    `gorm:"primary_key" json:"ID"`
	MenuType   string `validate:"required" json:"menu_type"`
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

// func (m Meal) UnMarshalJSON() ([]byte, error) {
// 	type meal Meal // prevent recursion
// 	x := meal(m)
// 	date = x.Date

// 	return json.Marshal(x)
// }
