package token

import (
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

func GenerateToken(user_id int) (Token, error) {
	secret_key := os.Getenv("JWT_SECRET_KEY")
	token := Token{}
	access_token_lifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MINUTE_LIFESPAN"))
	if err != nil {
		return token, err
	}
	refresh_token_lifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return token, err
	}
	access_token_claims := jwt.MapClaims{}
	access_token_claims["authorized"] = true
	access_token_claims["user_id"] = user_id
	access_token_claims["exp"] = time.Now().Add(time.Minute * time.Duration(access_token_lifespan)).Unix()
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, access_token_claims)

	refresh_token_claims := jwt.MapClaims{}
	refresh_token_claims["authorized"] = true
	refresh_token_claims["user_id"] = user_id
	refresh_token_claims["exp"] = time.Now().Add(time.Hour * time.Duration(refresh_token_lifespan)).Unix()

	access_token_string, err := access_token.SignedString([]byte(secret_key))
	if err != nil {
		return token, err
	}
	token.AccessToken = access_token_string

	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_token_claims)
	refresh_token_string, err := refresh_token.SignedString([]byte(secret_key))
	if err != nil {
		return token, err
	}
	token.RefreshToken = refresh_token_string

	return token, nil
}

func TokenValid(c *gin.Context) error {
	secret_key := os.Getenv("JWT_SECRET_KEY")
	claims := jwt.MapClaims{}
	token_string := ExtractToken(c)
	_, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
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
	secret_key := os.Getenv("JWT_SECRET_KEY")
	token_string := ExtractToken(c)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
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
