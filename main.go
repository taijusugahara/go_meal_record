package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "go_meal_record/app/db" //module名/ディレクトリ initだけの場合は_使用
	"net/http"
)

func main() {
	engine := gin.Default()
	//ログイン必要なし
	v1 := engine.Group("/v1")
	{
		account := v1.Group("/accounts")
		fmt.Println(account)

		{
			account.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "hello world",
				})
			})
			// account.GET("/list", controller.BookList)
			// account.POST("/register", controller.AccountRegister)
			// account.POST("/login", controller.Login)
		}
	}
	// 	//ログイン必要
	// 	v2 := bookEngine.Group("/v2")
	// 	v2.Use(middleware.JwtAuthMiddleware())
	// 	{
	// 		v2.GET("/user", controller.CurrentUser)
	// 		v2.DELETE("/user_delete", controller.UserDelete)
	// 		v2.POST("/add", controller.BookAdd)
	// 		v2.PUT("/update", controller.BookUpdate)
	// 		v2.DELETE("/delete", controller.BookDelete)
	// 	}
	engine.Static("/static", "./static")
	engine.Run(":3000")
}
