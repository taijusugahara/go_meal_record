package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
)

func MealIndexByDay(c *gin.Context) {
	meals := []model.Meal{}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tomorrow := today.AddDate(0, 0, 1)
	//BETWEENの場合today(00:00:00)含む。tomorrow(00:00:00)含む。
	result := db.DB.Where("meals.created_at BETWEEN ? AND ?", today, tomorrow).Preload("Menus").Preload("MealImages").Joins("User").Order("id").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
	}
	fmt.Println(time.Now().Format("2006-01-02 MST"))
	c.JSONP(http.StatusOK, gin.H{
		"data": meals,
	})
}

func MealIndexByMonth(c *gin.Context) {
	meals := []model.Meal{}
	now := time.Now()
	this_month := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local)
	next_month := this_month.AddDate(0, 1, 0)
	result := db.DB.Where("meals.created_at BETWEEN ? AND ?", this_month, next_month).Preload("Menus").Preload("MealImages").Joins("User").Order("id").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
	}
	fmt.Println(time.Now().Format("2006-01-02 MST"))
	c.JSONP(http.StatusOK, gin.H{
		"data": meals,
	})
}

func MenuIndex(c *gin.Context) {
	menus := []model.Menu{}
	result := db.DB.Order("id").Find(&menus)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
	}
	c.JSONP(http.StatusOK, gin.H{
		"data": menus,
	})
}

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
	//
	// now := time.Now()
	// today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// tomorrow := today.AddDate(0, 0, 1)
	//
	result := db.DB.Create(&meal)
	// meal.CreatedAt = today
	db.DB.Save(&meal)
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
	meal_id, err := strconv.Atoi(c.Param("meal_id"))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	meal := model.Meal{}
	menu := model.Menu{}
	err = c.BindJSON(&menu)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	result := db.DB.Joins("User").First(&meal, meal_id)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	menu.Meal = meal
	menu.MealID = meal_id
	err = validate.Struct(menu)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Validation Error")
		return
	}
	result = db.DB.Create(&menu)
	if result.Error != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"data": menu,
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
		"message": "delete success",
		"data":    menu,
	})

}

func MealImageCreate(c *gin.Context) {
	var err error
	meal_id, err := strconv.Atoi(c.Param("meal_id"))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	meal := model.Meal{}
	image := model.MealImage{}
	// err = c.Request.ParseMultipartForm(32 << 20)
	// if err != nil {
	// 	log.Println("c")
	// 	log.Println(err)
	// 	c.String(http.StatusInternalServerError, "Server Error")
	// 	return
	// }
	// err = c.Request.ParseMultipartForm(100000000)
	// if err != nil {
	// 	log.Println("b")
	// 	log.Println(err)
	// 	c.String(http.StatusInternalServerError, "Server Error")
	// 	return
	// }
	file, err := c.FormFile("file")
	// log.Println(file)
	if err != nil {
		log.Println("a")
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	result := db.DB.Joins("User").First(&meal, meal_id)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	if file != nil {
		filename := file.Filename
		filename_split_dot := strings.Split(filename, ".")
		extention := filename_split_dot[len(filename_split_dot)-1]
		log.Println(extention)
		valid_extentions := []string{"jpeg", "jpg", "JPEG", "png", "PNG"}
		if common.Contains(valid_extentions, extention) {
			image.File = filename
			image.Meal = meal
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
			path := image.File
			c.SaveUploadedFile(file, path)
		} else {
			c.String(http.StatusInternalServerError, "File extention not correct")
			return
		}

	}
	c.JSONP(http.StatusOK, gin.H{
		"file":  file,
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
		"message": "delete success",
		"data":    meal_image,
	})

}
