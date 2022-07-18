package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go_meal_record/app/controller"
	_ "go_meal_record/app/db" //module名/ディレクトリ initだけの場合は_使用
	"go_meal_record/app/middleware"

	"fmt"
	"net/http"
	"time"
)

func main() {
	engine := gin.Default()
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"http://localhost:3001"}
	cors_config.AllowHeaders = []string{"Access-Control-Allow-Credentials",
		"Access-Control-Allow-Headers",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization"}

	engine.Use(cors.New(cors_config))
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
		v2.POST("/get_token_by_refresh_token", controller.GetTokenByRefreshToken)
		v2.GET("/user_info", controller.ShowUser)

		meal := v2.Group("/meal")
		meal.GET("/index/day/", controller.MealIndexByDay)
		meal.GET("/index/month/", controller.MealIndexByMonth)
		meal.POST("/create/", controller.MealCreate)
		meal.GET("/menu_index/", controller.MenuIndex)
		meal.POST("/menu_create/:meal_id/", controller.MenuCreate)
		meal.POST("/image_create/:meal_id/", controller.MealImageCreate)
		// 		v2.DELETE("/user_delete", controller.UserDelete)
		// 		v2.POST("/add", controller.BookAdd)
		// 		v2.PUT("/update", controller.BookUpdate)
		// 		v2.DELETE("/delete", controller.BookDelete)
	}
	engine.Static("/static", "./static")
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	engine.Run(":3000")
}
