package test

import (
	"fmt"
	"go-meal-record/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJwtAuthMiddleware_SuccessWithToken(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	//handler(ただのfunc)の処理するのには下のようにする。
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	//req必要 urlはなんでもいい。
	//header付与
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", token.AccessToken))
	c.Request = request
	middleware.JwtAuthMiddleware(c)
	assert.Equal(t, 200, response.Code)
}
func TestJwtAuthMiddleware_FailWithoutToken(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	dummy_jwt_token_create(user.ID)
	//handler(ただのfunc)の処理するのには下のようにする。
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	//req必要 urlはなんでもいい。
	c.Request = request
	middleware.JwtAuthMiddleware(c)
	assert.Equal(t, 401, response.Code)
}
