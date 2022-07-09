package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go_meal_record/app/db"
	"go_meal_record/app/model"
	"go_meal_record/app/utils/token"

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
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	user.Password = string(hashedPassword)
	result := db.DB.Create(&user)
	if result.Error != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	//ログインさせる
	token, err := token.GenerateToken(user.ID)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "success!",
		"user":          user,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func Login(c *gin.Context) {
	var err error
	user := model.User{}
	email := c.PostForm("email")
	password := c.PostForm("password")
	user.Email = email
	user.Password = password

	my_user := model.User{}
	err = db.DB.Model(model.User{}).Where("email = ?", user.Email).First(&my_user).Error
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(my_user.Password), []byte(user.Password))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	token, err := token.GenerateToken(my_user.ID)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func GetTokenByRefreshToken(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	token, err := token.GenerateToken(user_id)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})

}

func GetCurrentUser(c *gin.Context) (model.User, error) {
	user := model.User{}
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		return user, err
	}
	err = db.DB.First(&user, user_id).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func ShowUser(c *gin.Context) {
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}

	log.Println(user)
	log.Println(user.Name)
	log.Println(user.Email)

	c.JSON(http.StatusOK, gin.H{"user": user})
}
