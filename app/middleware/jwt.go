package middleware

import (
	// "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"log"
	"net/http"

	"go-meal-record/app/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			log.Println(err)
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort() //以降の処理(main.go)しない
			return
		}
		c.Next() //main.goの処理に戻る
	}
}
