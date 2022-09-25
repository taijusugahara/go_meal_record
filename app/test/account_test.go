package test

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"bytes"
	"encoding/json"
	"go-meal-record/app/controller"
	"go-meal-record/app/db"
	"go-meal-record/app/model"
	"go-meal-record/app/utils/token"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutePath_Success(t *testing.T) {
	dbsetup()
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Message string
	}
	resp := ResponseData{}
	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, "Welcome to go-meal-record ok!!!", resp.Message)
}

func TestAccountRegister_Success(t *testing.T) {
	dbsetup()
	users := []model.User{}
	result := db.DB.Find(&users)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(users), 0)

	name := "AAA"
	email := "12@34.com"
	password := "kepelskeos"

	params := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}

	jsonParams, _ := json.Marshal(params)

	request, _ := http.NewRequest("POST", "/v1/account/register", bytes.NewBuffer(jsonParams))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	responseBody, _ := io.ReadAll(response.Body)

	type ResponseData struct {
		Message       string
		User          model.User
		Access_token  string
		Refresh_token string
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	//dbにuserが追加されているかどうか
	result = db.DB.Find(&users)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, 1, len(users))
	new_user := users[0]
	assert.Equal(t, name, new_user.Name)
	assert.Equal(t, email, new_user.Email)
	password_compare_err := bcrypt.CompareHashAndPassword([]byte(new_user.Password), []byte(password))
	assert.Equal(t, nil, password_compare_err)
	//

	if err != nil {
		log.Fatal(err)
		return
	}
	assert.Equal(t, "success!", resp.Message)
	assert.Equal(t, new_user.Name, resp.User.Name)
	//emailとpasswordはjsonで渡さない。
	assert.Equal(t, "", resp.User.Email)
	assert.Equal(t, "", resp.User.Password)
	//tokenはtest側で作成したものとcontroller側で作成したものが一致するわけではないので(userによって固定されるものでもない為)存在確認のみになる。
	assert.True(t, len(resp.Access_token) > 0)
	assert.True(t, len(resp.Refresh_token) > 0)
}

func TestAccountRegister_Fail(t *testing.T) {
	dbsetup()

	cases := []struct {
		name     string
		username string
		email    string
		password string
	}{
		{name: "without username", username: "", email: "12@34.com", password: "kepelskeos"},
		{name: "without email", username: "aiu", email: "", password: "kepelskeos"},
		{name: "not correct email", username: "aiu", email: "1234.com", password: "kepelskeos"},
		{name: "without password", username: "aiu", email: "12@34.com", password: ""},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			params := map[string]string{
				"name":     tt.username,
				"email":    tt.email,
				"password": tt.password,
			}
			jsonParams, _ := json.Marshal(params)
			request, _ := http.NewRequest("POST", "/v1/account/register", bytes.NewBuffer(jsonParams))
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)
			assert.Equal(t, 500, response.Code)
		})
	}
}

func TestLogin_Success(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	email := "xxx@yyy.com"
	password := "xxxyyyzzz"
	params := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonParams, _ := json.Marshal(params)
	request, _ := http.NewRequest("POST", "/v1/account/login", bytes.NewBuffer(jsonParams))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	responseBody, _ := io.ReadAll(response.Body)

	type ResponseData struct {
		Access_token  string
		Refresh_token string
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	token, _ := token.GenerateToken(user.ID)

	assert.Equal(t, token.AccessToken, resp.Access_token)
	assert.Equal(t, token.RefreshToken, resp.Refresh_token)
}

func TestLogin_Fail(t *testing.T) {
	dbsetup()
	dummy_user_create()

	cases := []struct {
		name     string
		email    string
		password string
	}{
		{name: "wrong email", email: "aaaaxxx@yyytest.com", password: "xxxyyyzzz"},
		{name: "wrong password", email: "xxx@yyytest.com", password: "aaaaaxxxyyyzzz"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			params := map[string]string{
				"email":    tt.email,
				"password": tt.password,
			}
			jsonParams, _ := json.Marshal(params)
			request, _ := http.NewRequest("POST", "/v1/account/login", bytes.NewBuffer(jsonParams))
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)
			assert.Equal(t, 500, response.Code)
		})
	}
}

func TestGetTokenByRefreshToken_Success(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)

	request, _ := http.NewRequest("POST", "/v2/get_token_by_refresh_token", nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", token.RefreshToken))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Access_token  string
		Refresh_token string
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}
	assert.True(t, len(resp.Access_token) > 0)
	assert.True(t, len(resp.Refresh_token) > 0)

}

func TestGetTokenByRefreshToken_FailWithoutToken(t *testing.T) {
	dbsetup()
	dummy_user_create()
	request, _ := http.NewRequest("POST", "/v2/get_token_by_refresh_token", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	assert.Equal(t, 401, response.Code)
}

func TestGetCurrentUser(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	//handler(ただのfunc)の処理するのには下のようにする。
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	//req必要 urlはなんでもいい。
	req, _ := http.NewRequest("GET", "/", nil)
	//header付与
	req.Header.Add("Authorization", fmt.Sprintf("JWT %v", token.AccessToken))
	c.Request = req
	log.Println(c)
	current_user, err := controller.GetCurrentUser(c)
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Name, current_user.Name)
}

func TestShowUser(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)

	request, _ := http.NewRequest("GET", "/v2/user_info", nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", token.AccessToken))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		User model.User
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, user.Name, resp.User.Name)
	//emailとpasswordは渡さない
	assert.Equal(t, "", resp.User.Email)
	assert.Equal(t, "", resp.User.Password)
}
