package main

import (
	// "meal_record_practice/app/controller"
	// "meal_record_practice/app/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	// bookEngine := engine.Group("/book")
	// {
	// 	//ログイン必要なし
	// 	v1 := bookEngine.Group("/v1")
	// 	{
	// 		v1.GET("/list", controller.BookList)
	// 		v1.POST("/register", controller.Register)
	// 		v1.POST("/login", controller.Login)
	// 	}
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
	// }
	engine.Static("/static", "./static")
	engine.Run(":3000")
}
