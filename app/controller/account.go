package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-meal-record/app/db"
	"go-meal-record/app/model"
	"go-meal-record/app/utils/token"
	"go-meal-record/app/utils/validate"

	"log"
	"net/http"
)

//log.Fatallnだとエラー起きるとサーバーが落ちる

/////c.PostForm()はpostmanのformdataの時にのみ動いた。
//postmanのrawやreactからはjsonで渡ってくるのでBindJSONを使う感じ。
//BindJSONを使っておこう。apiの挙動確認もrawでやったほうが良さそう

func AccountRegister(c *gin.Context) {
	var err error
	validate := validate.Validate()
	user := model.User{}
	err = c.BindJSON(&user)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	err = validate.Struct(user)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Validation Error")
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
	//Nameないのに&userでも大丈夫な理由はvalidateしてないから。validateする場合はnameなしのstruct作ること
	err = c.BindJSON(&user)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

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
	c.JSON(http.StatusOK, gin.H{"user": user})
}
