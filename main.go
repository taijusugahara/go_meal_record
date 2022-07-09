package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_meal_record/app/controller"
	_ "go_meal_record/app/db" //module名/ディレクトリ initだけの場合は_使用
	"go_meal_record/app/middleware"
	"net/http"
)

func main() {
	engine := gin.Default()
	//ログイン必要なし
	v1 := engine.Group("/v1")
	{
		account := v1.Group("/account")
		fmt.Println(account)

		{
			// account.GET("/list", controller.BookList)
			account.POST("/register", controller.AccountRegister)
			account.POST("/login", controller.Login)
		}
	}
	// 	//ログイン必要
	v2 := engine.Group("/v2")
	v2.Use(middleware.JwtAuthMiddleware())
	{
		v2.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello world",
			})
		})
		v2.GET("/user", controller.ShowUser)
		// 		v2.DELETE("/user_delete", controller.UserDelete)
		// 		v2.POST("/add", controller.BookAdd)
		// 		v2.PUT("/update", controller.BookUpdate)
		// 		v2.DELETE("/delete", controller.BookDelete)
	}
	engine.Static("/static", "./static")
	engine.Run(":3000")
}
