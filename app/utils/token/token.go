package token

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(user_id int) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("token"))

}

func TokenValid(c *gin.Context) error {
	token_string := ExtractToken(c)
	claims := jwt.MapClaims{}
	if token_string == "" {
		log.Println("token is not exist")
		c.String(http.StatusUnauthorized, "Unauthorized")
	}
	_, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("token"), nil
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(c *gin.Context) (int, error) {

	token_string := ExtractToken(c)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("token"), nil
	})
	if err != nil {
		return 0, err
	}
	claims_user_id := claims["user_id"]
	//interface型でuser_idはfloat型らしい。なので１回floatで取り出してからintにしてる
	//type assertionしてからintにしてる
	user_id := int(claims_user_id.(float64))
	return user_id, nil
}
