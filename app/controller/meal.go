package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_meal_record/app/db"
	"go_meal_record/app/model"

	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	validate := validator.New()
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
		c.String(http.StatusInternalServerError, "Server Error")
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
	validate := validator.New()
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
		c.String(http.StatusInternalServerError, "Server Error")
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
	file, err := c.FormFile("file")
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

	if file != nil {
		filename := file.Filename
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

	}
	c.JSONP(http.StatusOK, gin.H{
		"file":  file,
		"image": image,
	})
}
