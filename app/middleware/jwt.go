package middleware

import (
	"github.com/gin-gonic/gin"

	"log"
	"net/http"

	"go-meal-record/app/utils/token"
)

func JwtAuthMiddleware(c *gin.Context) {
	err := token.TokenValid(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort() //以降の処理(main.go)しない
		return
	}
	c.Next() //main.goの処理に戻る
}
