package model

import (
	// "log"
	// "os"
	// "strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"go_meal_record/app/utils/common"
)

type Meal struct {
	gorm.Model `json:"-"`
	ID         int            `gorm:"primary_key" json:"ID"`
	MealType   string         `validate:"meal_type,required" json:"meal_type"` //custom validate
	User       User           `json:"user" validate:"required"`
	UserID     int            `json:"user_id validate:"required"`
	Menus      []Menu         `json:"menus"`
	MealImages []MealImage    `json:"meal_images"`
	Date       datatypes.Date `json:"date" validate:"required"`
} //datatypes.Dateは日時の部分を00:00:00にしてくれる

type Menu struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primary_key" json:"ID"`
	Name       string `validate:"required,max=255" json:"name"`
	Meal       Meal   `json:"-" validate:"required"`
	MealID     int    `json:"meal_id" validate:"required"`
}

type MealImage struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primary_key" json:"ID"`
	Filename   string `json:"-"`
	Fileurl    string `json:"fileurl"`
	Meal       Meal   `json:"-" validate:"required"`
	MealID     int    `json:"meal_id" validate:"required"`
}

//#
// jsonで返すときに外部キーの情報も全て渡ってしまうので無駄な情報が多い。
// 例でいうとMealの中のMealImageの中にMenuの情報が入っている。
// 別のstructを作って必要なカラムだけ入れようとしたが失敗。many_fieldに外部キー指定されているのが原因。
// json:"-"とすることでjsonで情報が渡らないようにしてる。これはpostでもjsonを送れなくなるのでpost必要なカラムは無理。
// 一覧表示などの場合はMealImageの中のMenu情報はいらないと思うが、
// 個の情報を取るときはMenuが欲しい場合もあるかもしれない。
// けどReactであればMealImageとMenuの情報を両方渡せばいいし、セットになってなくてもいいと思う。mapで展開して表示してるからデータ自体は持ってる
//#

//Meal{}を作成したところデフォルトのMealImagesとMenusが[]ではなくnilで返ってくるので[]にしてあげる
func NewMeal() Meal {
	m := Meal{}
	m.MealImages = make([]MealImage, 0)
	m.Menus = make([]Menu, 0)
	return m
}

func CustomValidateMealType(fl validator.FieldLevel) bool {
	meal_type_choices := []string{"morning", "lunch", "dinner", "other"}
	if common.Contains(meal_type_choices, fl.Field().String()) {
		return true
	}
	return false
}
