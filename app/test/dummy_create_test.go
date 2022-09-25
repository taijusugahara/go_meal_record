package test

import (
	"bytes"
	"fmt"
	"go-meal-record/app/db"
	"go-meal-record/app/model"
	"go-meal-record/app/utils/token"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
)

//なぜか分からないが.._test.goにしないとcommon_test.goで共有されてるengineの部分がこちらで赤色エラー表示される。testは通る
//その為このfileではtestを行なっていないのだが、_test.goの形にしてる

func dummy_user_create() model.User {
	user := model.User{}
	user.Name = "test"
	user.Email = "xxx@yyy.com"
	password := "xxxyyyzzz"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	db.DB.Create(&user)
	return user
}

func dummy_user_create2() model.User {
	user := model.User{}
	user.Name = "test2"
	user.Email = "aaa@bbb.com"
	password := "xxxyyyzzz"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	db.DB.Create(&user)
	return user
}

func dummy_jwt_token_create(user_id int) token.Token {
	token, _ := token.GenerateToken(user_id)
	return token
}

func dummy_meal_create(access_token_list []string, date_str_list []string) {
	menu_name := "tekito"
	meal_type_list := []string{"morning", "lunch", "dinner", "other"}

	// user二人分作成
	for _, access_token := range access_token_list {
		for _, date_str := range date_str_list {
			for _, meal_type := range meal_type_list {
				// formdata test 参考 https://stackoverflow.com/questions/59889551/testing-a-multipart-form-upload-endpoint-in-echo-framework

				//menu作成
				bodyBuf := &bytes.Buffer{}
				bodyWriter := multipart.NewWriter(bodyBuf)
				bodyWriter.WriteField("meal_type", meal_type)
				bodyWriter.WriteField("date", date_str)
				bodyWriter.WriteField("menu", menu_name)
				bodyWriter.Close() // <<< important part
				request, _ := http.NewRequest("POST", "/v2/meal/menu_create/", bodyBuf)
				request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
				request.Header.Set("Content-Type", bodyWriter.FormDataContentType()) // <<< important part
				response := httptest.NewRecorder()
				engine.ServeHTTP(response, request)

				//mealimage作成
				file, _ := os.Open("demo1.png") //io.copyは1回してしまうと、2回目以降はfile空になったので、毎度os.openしてる
				defer file.Close()
				bodyBuf2 := &bytes.Buffer{}
				bodyWriter2 := multipart.NewWriter(bodyBuf2)
				bodyWriter2.WriteField("meal_type", meal_type)
				bodyWriter2.WriteField("date", date_str)
				part, _ := bodyWriter2.CreateFormFile("file", "demo1.png")
				io.Copy(part, file)
				bodyWriter2.Close() // <<< important part
				request2, _ := http.NewRequest("POST", "/v2/meal/image_create", bodyBuf2)
				request2.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
				request2.Header.Set("Content-Type", bodyWriter2.FormDataContentType()) // <<< important part
				response2 := httptest.NewRecorder()
				engine.ServeHTTP(response2, request2)
			}
		}
	}
}
