package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go_meal_record/app/db"
	"go_meal_record/app/model"
	"net/http"
	// "strconv"

	"log"
)

//log.Fatallnだとエラー起きるとサーバーが落ちる

func AccountRegister(c *gin.Context) {
	var err error
	validate := validator.New()
	user := model.User{}
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	user.Name = name
	user.Email = email
	user.Password = password

	err = validate.Struct(user)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}
	user.Password = string(hashedPassword)
	result := db.DB.Create(&user)
	if result.Error != nil {
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success!"})
}
