package controller

import (
	"fmt"
	"go_meal_record/app/db"
	"go_meal_record/app/model"
	"go_meal_record/app/utils/common"
	"go_meal_record/app/utils/validate"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func MealIndexByDay(c *gin.Context) {
	date_str := c.Param("date")
	//2006-01-.....これじゃないといけないみたい
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	//dateはmodelでdatatypes.Date使用しているので例えば2022:07:20T10:40:30Zは2022:07:20T:00:00:00Zとなる。
	//そのため時間は00:00:00,0固定でいい。
	meals := []model.Meal{}
	one_day_meal_map := map[string]interface{}{
		"morning": model.NewMeal(), "lunch": model.NewMeal(), "dinner": model.NewMeal(), "other": model.NewMeal(),
	} //注意 : mapは強制的にkey abc順になる、morningから開始したいがslice rangeで順番変えても無理だった。今回の場合はdinnerが先頭に来る。morningからにしたかったらjs側でsortするか

	result := db.DB.Model(model.Meal{}).Where("date = ?", date_start).Preload("Menus").Preload("MealImages").Joins("User").Order("id").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	for _, meal := range meals {
		meal_type := meal.MealType
		one_day_meal_map[meal_type] = meal
	}

	c.JSONP(http.StatusOK, gin.H{
		"data": one_day_meal_map,
	})
}

func MealIndexByWeek(c *gin.Context) {
	meals := []model.Meal{}
	date_str := c.Param("date")
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	weekday := date.Weekday()
	weekday_int := int(weekday)
	first_week_date := date_start.AddDate(0, 0, -weekday_int)
	end_week_date := date_start.AddDate(0, 0, 6-weekday_int)
	result := db.DB.Where("meals.date BETWEEN ? AND ?", first_week_date, end_week_date).Preload("Menus").Preload("MealImages").Joins("User").Order("id").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"data": meals,
	})
}

func MealIndexByMonth(c *gin.Context) {
	meals := []model.Meal{}
	date_str := c.Param("date")
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	this_month := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	next_month := this_month.AddDate(0, 1, -1)
	log.Println(next_month)
	result := db.DB.Where("meals.date BETWEEN ? AND ?", this_month, next_month).Preload("Menus").Preload("MealImages").Joins("User").Order("id").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"data": meals,
	})
}

// func MenuIndex(c *gin.Context) {
// 	menus := []model.Menu{}
// 	result := db.DB.Order("id").Find(&menus)
// 	if result.Error != nil {
// 		log.Println(result.Error)
// 		c.String(http.StatusInternalServerError, "Server Error")
// 	}
// 	c.JSONP(http.StatusOK, gin.H{
// 		"data": menus,
// 	})
// }

func MealCreate(c *gin.Context) {
	validate := validate.Validate()
	var err error
	meal := model.Meal{}
	err = c.BindJSON(&meal)
	fmt.Println(err)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	meal.User = user
	meal.UserID = user.ID
	err = validate.Struct(meal)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Validation Error")
		return
	}
	date := meal.Date
	previous_meal := model.Meal{}
	//dateはmodelでdatatypes.Date使用しているので例えば2022:07:20T10:40:30Zは2022:07:20T:00:00:00Zとなる
	is_already_create := db.DB.Where("meal_type = ? AND user_id = ? AND date = ?", meal.MealType, meal.UserID, date).First(&previous_meal)
	//既に作成している場合は作成させない。
	if is_already_create.Error == nil {
		c.String(http.StatusInternalServerError, "The record already exist so cannot create.")
		return
	}
	//
	// now := time.Now()
	// today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// tomorrow := today.AddDate(0, 0, 1)
	//
	result := db.DB.Create(&meal)
	// meal.CreatedAt = today
	// db.DB.Save(&meal)
	if result.Error != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	log.Println(meal.ID)

	c.JSONP(http.StatusOK, gin.H{
		"data": meal,
	})
}

func MenuCreate(c *gin.Context) {
	validate := validate.Validate()
	var err error
	meal := model.Meal{}
	menu := model.Menu{}
	//js FormDataから渡してる。JsonBindしないのは渡ってきた3つが1つのstructに入るわけではないから。
	meal_type := c.PostForm("meal_type")
	date_str := c.PostForm("date")
	menu_name := c.PostForm("menu")

	already_created_meal := model.Meal{}
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := datatypes.Date(date)
	// ErrRecordNotFound エラーを避けたい場合は、db.Limit(1).Find(&user)のように、Find を 使用することができます。
	db.DB.Where("meal_type = ? AND user_id = ? AND date = ?", meal_type, user_id, date_start).Joins("User").Limit(1).Find(&already_created_meal)
	if already_created_meal.ID != 0 {
		meal = already_created_meal
	} else { //新しく作成する
		meal.MealType = meal_type
		meal.Date = date_start
		meal.User = user
		meal.UserID = user_id

		err = validate.Struct(meal)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Validation Error")
			return
		}

		result := db.DB.Create(&meal)
		if result.Error != nil {
			log.Println(result.Error)
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}
	}
	//同じ名前とmeal_idを持つmenuは作成させない。
	already_created_menu := model.Menu{}
	is_already_menu_result := db.DB.Where("name = ? AND meal_id = ?", menu_name, meal.ID).First(&already_created_menu)
	if is_already_menu_result.Error == nil { //存在したら
		log.Println(is_already_menu_result.Error)
		c.String(http.StatusInternalServerError, "既に同じ名前とmeal_idを持つmenuが存在するため作成できません。")
		return
	}

	menu.Name = menu_name
	menu.Meal = meal
	menu.MealID = meal.ID
	err = validate.Struct(menu)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Validation Error")
		return
	}
	result := db.DB.Create(&menu)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"menu": menu,
	})
}

func MenuDelete(c *gin.Context) {
	id := c.Param("id") //dbにidを渡す際、stringでもintでもどっちもでもいいみたい。
	menu := model.Menu{}
	result := db.DB.Joins("Meal").First(&menu, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID
	if menu.Meal.UserID != user_id {
		c.String(http.StatusBadRequest, "User not correct user")
		return
	}
	result = db.DB.Delete(&menu)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"menu": menu,
	})

}

func MealImageCreate(c *gin.Context) {
	var err error
	validate := validate.Validate()
	meal := model.Meal{}
	image := model.MealImage{}
	//fileがある場合はjsonBindじゃなくてc.PostForm使っていけた。
	//理由としてはheaderのcontent-typeをmultipart-form;にしてるからだと思う。いやデータ格納するのにFormData使ってるからかな
	meal_type := c.PostForm("meal_type")
	date_str := c.PostForm("date")
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	//Menuがあるか確認。なければ作成
	already_created_meal := model.Meal{}
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := datatypes.Date(date)
	// ErrRecordNotFound エラーを避けたい場合は、db.Limit(1).Find(&user)のように、Find を 使用することができます。
	db.DB.Where("meal_type = ? AND user_id = ? AND date = ?", meal_type, user_id, date_start).Joins("User").Limit(1).Find(&already_created_meal)
	//既に作成している場合
	if already_created_meal.ID != 0 {
		meal = already_created_meal
	} else { //新しく作成する
		meal.MealType = meal_type
		meal.Date = date_start
		meal.User = user
		meal.UserID = user_id

		err = validate.Struct(meal)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Validation Error")
			return
		}

		result := db.DB.Create(&meal)
		if result.Error != nil {
			log.Println(result.Error)
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}
	}

	if file != nil {
		filename := file.Filename
		filename_split_dot := strings.Split(filename, ".")
		extention := filename_split_dot[len(filename_split_dot)-1]
		valid_extentions := []string{"jpeg", "jpg", "JPEG", "png", "PNG"}
		if common.Contains(valid_extentions, extention) {
			image.File = filename
			image.Meal = meal
			image.MealID = meal.ID
			err = validate.Struct(image)
			if err != nil {
				log.Println(err)
				c.String(http.StatusInternalServerError, "Validation Error")
				return
			}
			result := db.DB.Create(&image)
			if result.Error != nil {
				log.Println(result.Error)
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}
			id := strconv.Itoa(image.ID)
			err := os.Mkdir("app/static/meal/"+id, 0750)
			if err != nil {
				log.Println(err)
				c.String(http.StatusInternalServerError, "Server Error")
				return
			}
			path := "app/" + image.File
			c.SaveUploadedFile(file, path)
		} else {
			c.String(http.StatusInternalServerError, "File extention not correct")
			return
		}

	}
	c.JSONP(http.StatusOK, gin.H{
		"image": image,
	})
}

func MealImageDelete(c *gin.Context) {
	id := c.Param("id") //dbにidを渡す際、stringでもintでもどっちもでもいいみたい。
	meal_image := model.MealImage{}
	result := db.DB.Joins("Meal").First(&meal_image, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID
	if meal_image.Meal.UserID != user_id {
		c.String(http.StatusBadRequest, "User not correct user")
		return
	}
	result = db.DB.Delete(&meal_image)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"meal_image": meal_image,
	})

}
