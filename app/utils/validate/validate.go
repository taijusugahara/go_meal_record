package validate

import (
	"github.com/go-playground/validator/v10"

	"go-meal-record/app/model"
)

func Validate() *validator.Validate {
	validate := validator.New()
	//custom validate add
	//外部キーにcustom validateがある場合でもaddしてないとエラー出るので共通化
	validate.RegisterValidation("meal_type", model.CustomValidateMealType)
	return validate
}
