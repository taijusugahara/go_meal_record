package main

import (
	"github.com/gin-contrib/cors"
	// limits "github.com/gin-contrib/size"
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
	// engine.MaxMultipartMemory = 1000 << 20
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"http://localhost:3001"}
	cors_config.AllowHeaders = []string{"Access-Control-Allow-Credentials",
		"Access-Control-Allow-Headers",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization"}

	engine.Use(cors.New(cors_config))
	// engine.Use(limits.RequestSizeLimiter(100000000))
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
		meal.GET("/index/day/:date/", controller.MealIndexByDay)
		meal.GET("/index/week/:date/", controller.MealIndexByWeek)
		meal.GET("/index/month/:date/", controller.MealIndexByMonth)
		meal.POST("/create/", controller.MealCreate)
		// meal.GET("/menu_index/", controller.MenuIndex)
		meal.POST("/menu_create/", controller.MenuCreate)
		//!!!!!!!下、最後に/(スラッシュ)つけてはいけない。/つけた場合は1MB以上のファイルを送信できなかった(POSTMAN)。Reactは試してない。https://stackoverflow.com/questions/33771167/handle-file-uploading-with-go
		meal.POST("/image_create", controller.MealImageCreate)
		meal.DELETE("/delete/menu/:id/", controller.MenuDelete)
		meal.DELETE("/delete/meal_image/:id/", controller.MealImageDelete)
	}
	engine.Static("/static", "./app/static")
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	engine.Run(":3000")
}
