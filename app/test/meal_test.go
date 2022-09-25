package test

import (
	"fmt"
	"go-meal-record/app/db"
	"go-meal-record/app/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"encoding/json"

	"github.com/stretchr/testify/assert"

	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMealIndexByDay(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	another_user := dummy_user_create2()
	another_token := dummy_jwt_token_create(another_user.ID)
	another_access_token := another_token.AccessToken
	access_token_list := []string{access_token, another_access_token}

	meal_type_list := []string{"morning", "lunch", "dinner", "other"}
	//対象日時の24日とその前後の23,25のmeal作成する。24にアクセスするので返ってくるのは24だけ。
	date_str_list := []string{"2022-07-23T00:00:00Z", "2022-07-24T00:00:00Z", "2022-07-25T00:00:00Z"}
	date_str := "2022-07-24T00:00:00Z"
	dummy_meal_create(access_token_list, date_str_list)
	request, _ := http.NewRequest("GET", fmt.Sprintf("/v2/meal/index/day/%v/", date_str), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Data map[string]interface{}
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	for _, meal_type := range meal_type_list {
		//interfaceから値取り出すには以下のようにする(型アサーション)
		resp_meal_type_data := resp.Data[meal_type].(map[string]interface{})

		meal_images := resp_meal_type_data["meal_images"].([]interface{})
		menus := resp_meal_type_data["menus"].([]interface{})
		meal_user := resp_meal_type_data["user"].(map[string]interface{})

		assert.Equal(t, date_str, resp_meal_type_data["date"])
		assert.Equal(t, 1, len(meal_images))
		assert.Equal(t, 1, len(menus))
		assert.Equal(t, user.ID, int(meal_user["ID"].(float64)))
		assert.Equal(t, "test", meal_user["name"])
	}
}

func TestMealIndexByWeek(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	another_user := dummy_user_create2()
	another_token := dummy_jwt_token_create(another_user.ID)
	another_access_token := another_token.AccessToken
	access_token_list := []string{access_token, another_access_token}

	meal_type_list := []string{"morning", "lunch", "dinner", "other"}
	//対象日時の2022/07/24日の週は[24,25,26,27,28,29,30]である。23と31は範囲外要素として検証
	date_str_list := []string{"2022-07-23T00:00:00Z", "2022-07-24T00:00:00Z", "2022-07-25T00:00:00Z", "2022-07-26T00:00:00Z", "2022-07-27T00:00:00Z", "2022-07-28T00:00:00Z", "2022-07-29T00:00:00Z", "2022-07-30T00:00:00Z", "2022-07-31T00:00:00Z"}
	date_str := "2022-07-24T00:00:00Z"
	dummy_meal_create(access_token_list, date_str_list)
	request, _ := http.NewRequest("GET", fmt.Sprintf("/v2/meal/index/week/%v/", date_str), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Data            map[string]map[string]interface{}
		First_week_date string
		End_week_date   string
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, "2022-07-24T00:00:00Z", resp.First_week_date)
	assert.Equal(t, "2022-07-30T00:00:00Z", resp.End_week_date)

	for _, date_str := range date_str_list {
		for _, meal_type := range meal_type_list {
			if date_str == "2022-07-23T00:00:00Z" || date_str == "2022-07-31T00:00:00Z" {
				_, ok := resp.Data[date_str][meal_type]
				assert.False(t, ok)
			} else {
				resp_meal_type_data := resp.Data[date_str][meal_type].(map[string]interface{})
				meal_images := resp_meal_type_data["meal_images"].([]interface{})
				menus := resp_meal_type_data["menus"].([]interface{})
				meal_user := resp_meal_type_data["user"].(map[string]interface{})

				assert.Equal(t, date_str, resp_meal_type_data["date"])
				assert.Equal(t, 1, len(meal_images))
				assert.Equal(t, 1, len(menus))
				assert.Equal(t, user.ID, int(meal_user["ID"].(float64)))
				assert.Equal(t, "test", meal_user["name"])
			}
		}
	}
}

func TestMealIndexByMonth(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	another_user := dummy_user_create2()
	another_token := dummy_jwt_token_create(another_user.ID)
	another_access_token := another_token.AccessToken
	access_token_list := []string{access_token, another_access_token}

	meal_type_list := []string{"morning", "lunch", "dinner", "other"}
	//対象日時の2022/07/の月は[1,24,31]である。06/30と08/01は範囲外要素として検証
	date_str_list := []string{"2022-06-30T00:00:00Z", "2022-07-01T00:00:00Z", "2022-07-24T00:00:00Z", "2022-07-31T00:00:00Z", "2022-08-01T00:00:00Z"}
	date_str := "2022-07-24T00:00:00Z"
	dummy_meal_create(access_token_list, date_str_list)
	request, _ := http.NewRequest("GET", fmt.Sprintf("/v2/meal/index/month/%v/", date_str), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Data             map[string]map[string]interface{}
		First_month_date string
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, "2022-07-01T00:00:00Z", resp.First_month_date)

	for _, date_str := range date_str_list {
		for _, meal_type := range meal_type_list {
			if date_str == "2022-06-30T00:00:00Z" || date_str == "2022-08-01T00:00:00Z" {
				_, ok := resp.Data[date_str][meal_type]
				assert.False(t, ok)
			} else {
				resp_meal_type_data := resp.Data[date_str][meal_type].(map[string]interface{})
				meal_images := resp_meal_type_data["meal_images"].([]interface{})
				menus := resp_meal_type_data["menus"].([]interface{})
				meal_user := resp_meal_type_data["user"].(map[string]interface{})

				assert.Equal(t, date_str, resp_meal_type_data["date"])
				assert.Equal(t, 1, len(meal_images))
				assert.Equal(t, 1, len(menus))
				assert.Equal(t, user.ID, int(meal_user["ID"].(float64)))
				assert.Equal(t, "test", meal_user["name"])
			}
		}
	}
}

func TestMenuCreate_SuccessWithoutMeal(t *testing.T) {
	//mealがない場合mealを作成してmenuを作成する。
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken

	meals := []model.Meal{}
	result := db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 0)

	menus := []model.Menu{}
	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 0)

	menu_name := "apple"
	meal_type := "morning"
	date_str := "2022-07-24T00:00:00Z"
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

	assert.Equal(t, 200, response.Code)

	//dbに保存されてるかどうか

	result = db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	meal := meals[0]
	date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
	meal_date_str := date.Format("2006-01-02T15:04:05Z")
	assert.Equal(t, date_str, meal_date_str)
	assert.Equal(t, meal.UserID, user.ID)

	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 1)

	menu := menus[0]
	assert.Equal(t, menu_name, menu.Name)
	assert.Equal(t, meal.ID, menu.MealID)

	//

	//jsonの中身

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Menu model.Menu
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, menu_name, resp.Menu.Name)
	assert.Equal(t, meal.ID, resp.Menu.MealID)

}

func TestMenuCreate_SuccessWithMeal(t *testing.T) {

	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	//1回目 mealなし
	menu_name := "banana"
	meal_type := "morning"
	date_str := "2022-07-24T00:00:00Z"
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

	meals := []model.Meal{}
	result := db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	menus := []model.Menu{}
	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 1)

	//2回目 mealあり

	menu_name2 := "apple"
	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	bodyWriter2.WriteField("meal_type", meal_type)
	bodyWriter2.WriteField("date", date_str)
	bodyWriter2.WriteField("menu", menu_name2)
	bodyWriter2.Close() // <<< important part
	request2, _ := http.NewRequest("POST", "/v2/meal/menu_create/", bodyBuf2)
	request2.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	request2.Header.Set("Content-Type", bodyWriter2.FormDataContentType()) // <<< important part
	response2 := httptest.NewRecorder()
	engine.ServeHTTP(response2, request2)

	assert.Equal(t, 200, response2.Code)

	//dbに保存されてるかどうか

	result = db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	meal := meals[0]
	date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
	meal_date_str := date.Format("2006-01-02T15:04:05Z")
	assert.Equal(t, date_str, meal_date_str)
	assert.Equal(t, meal.UserID, user.ID)

	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 2)
	menu2 := menus[1]
	assert.Equal(t, menu_name2, menu2.Name)
	assert.Equal(t, meal.ID, menu2.MealID)

	//

	//jsonの中身

	responseBody, _ := io.ReadAll(response2.Body)
	type ResponseData struct {
		Menu model.Menu
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, menu_name2, resp.Menu.Name)
	assert.Equal(t, meal.ID, resp.Menu.MealID)

}

func TestMenuCreate_FailSameMenuInSameMeal(t *testing.T) {
	//同じmealないであれば同じ名前のmenuを作成できない。
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	//1回目 mealなし
	menu_name := "banana"
	meal_type := "morning"
	date_str := "2022-07-24T00:00:00Z"
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

	meals := []model.Meal{}
	result := db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	menus := []model.Menu{}
	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 1)

	//2回目 mealあり

	menu_name2 := "banana"
	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	bodyWriter2.WriteField("meal_type", meal_type)
	bodyWriter2.WriteField("date", date_str)
	bodyWriter2.WriteField("menu", menu_name2)
	bodyWriter2.Close() // <<< important part
	request2, _ := http.NewRequest("POST", "/v2/meal/menu_create/", bodyBuf2)
	request2.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	request2.Header.Set("Content-Type", bodyWriter2.FormDataContentType()) // <<< important part
	response2 := httptest.NewRecorder()
	engine.ServeHTTP(response2, request2)

	assert.Equal(t, 500, response2.Code)

	//dbに保存されてるかどうか

	result = db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	meal := meals[0]
	date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
	meal_date_str := date.Format("2006-01-02T15:04:05Z")
	assert.Equal(t, date_str, meal_date_str)
	assert.Equal(t, meal.UserID, user.ID)

	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 1)
	menu := menus[0]
	assert.Equal(t, menu_name, menu.Name)
	assert.Equal(t, meal.ID, menu.MealID)

}

func TestMenuDelete_Success(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	//1度 meal,menu作成
	access_token_list := []string{access_token}
	date_str_list := []string{"2022-07-24T00:00:00Z"}
	dummy_meal_create(access_token_list, date_str_list)

	menus := []model.Menu{}
	result := db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 4)

	menu := menus[0]
	menu_id := menu.ID

	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/v2/meal/delete/menu/%v/", menu_id), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)

	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 3)
	for _, m := range menus {
		assert.NotEqual(t, menu_id, m.ID)
	}

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Menu model.Menu
	}

	resp := ResponseData{}

	err := json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, menu_id, resp.Menu.ID)
}

func TestMenuDelete_FailDifferentUser(t *testing.T) {
	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	user2 := dummy_user_create2()
	token2 := dummy_jwt_token_create(user2.ID)
	access_token2 := token2.AccessToken
	//1度 meal作成
	access_token_list := []string{access_token}
	date_str_list := []string{"2022-07-24T00:00:00Z"}
	dummy_meal_create(access_token_list, date_str_list)

	menus := []model.Menu{}
	result := db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 4)

	menu := menus[0]
	menu_id := menu.ID

	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/v2/meal/delete/menu/%v/", menu_id), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token2))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code)

	result = db.DB.Find(&menus)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menus), 4)

}

func TestMealImageCreate_SuccessWithoutMeal(t *testing.T) {
	//mealなしの場合はmealを作成してからmeal_imageを作成する
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsS3Bucket := os.Getenv("AWS_S3_TEST_BUCKET")

	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken

	meals := []model.Meal{}
	result := db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 0)

	meal_images := []model.MealImage{}
	result = db.DB.Find(&meal_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meal_images), 0)

	file, _ := os.Open("demo1.png")
	defer file.Close()
	meal_type := "morning"
	date_str := "2022-07-24T00:00:00Z"
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("meal_type", meal_type)
	bodyWriter.WriteField("date", date_str)
	part, _ := bodyWriter.CreateFormFile("file", "demo1.png")
	io.Copy(part, file)
	bodyWriter.Close() // <<< important part
	request, _ := http.NewRequest("POST", "/v2/meal/image_create", bodyBuf)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType()) // <<< important part
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	//dbに保存されているかどうか
	result = db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	meal := meals[0]
	date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
	meal_date_str := date.Format("2006-01-02T15:04:05Z")
	assert.Equal(t, date_str, meal_date_str)
	assert.Equal(t, meal.UserID, user.ID)

	result = db.DB.Find(&meal_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meal_images), 1)

	meal_image := meal_images[0]

	image_id := strconv.Itoa(meal_image.ID)
	filename := "meal/" + image_id + "/" + "demo1.png"
	file_url := "https://" + awsS3Bucket + ".s3." + awsRegion + ".amazonaws.com/" + filename
	assert.Equal(t, meal.ID, meal_image.MealID)
	assert.Equal(t, filename, meal_image.Filename)
	assert.Equal(t, file_url, meal_image.Fileurl)

	//

	//s3にあるか確認
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey, awsSecretKey, "",
		),
	}))

	download_file, _ := os.Create("downloader.png")
	defer download_file.Close()

	downloader := s3manager.NewDownloader(newSession)
	numBytes, err := downloader.Download(download_file,
		&s3.GetObjectInput{
			Bucket: aws.String(awsS3Bucket),
			Key:    aws.String(filename),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, true, numBytes > 0)

	os.Remove("downloader.png") //削除しておく

	//

	// json確認

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Image model.MealImage
	}

	resp := ResponseData{}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, meal_image.ID, resp.Image.ID)
	assert.Equal(t, file_url, resp.Image.Fileurl)
}

func TestMealImageCreate_SuccessWithMeal(t *testing.T) {
	//mealなしの場合はmealを作成してからmeal_imageを作成する
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsS3Bucket := os.Getenv("AWS_S3_TEST_BUCKET")

	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken

	//1回目 mealなし
	file, _ := os.Open("demo1.png")
	defer file.Close()
	meal_type := "morning"
	date_str := "2022-07-24T00:00:00Z"
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("meal_type", meal_type)
	bodyWriter.WriteField("date", date_str)
	part, _ := bodyWriter.CreateFormFile("file", "demo1.png")
	io.Copy(part, file)
	bodyWriter.Close() // <<< important part
	request, _ := http.NewRequest("POST", "/v2/meal/image_create", bodyBuf)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType()) // <<< important part
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	meals := []model.Meal{}
	result := db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	meal_images := []model.MealImage{}
	result = db.DB.Find(&meal_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meal_images), 1)

	//2回目 mealあり
	file2, _ := os.Open("demo1.png")
	defer file2.Close()
	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	bodyWriter2.WriteField("meal_type", meal_type)
	bodyWriter2.WriteField("date", date_str)
	part2, _ := bodyWriter2.CreateFormFile("file", "demo1.png")
	io.Copy(part2, file2)
	bodyWriter2.Close() // <<< important part
	request2, _ := http.NewRequest("POST", "/v2/meal/image_create", bodyBuf2)
	request2.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	request2.Header.Set("Content-Type", bodyWriter2.FormDataContentType()) // <<< important part
	response2 := httptest.NewRecorder()
	engine.ServeHTTP(response2, request2)

	//dbに保存されているかどうか
	result = db.DB.Find(&meals)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meals), 1)

	meal := meals[0]
	date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
	meal_date_str := date.Format("2006-01-02T15:04:05Z")
	assert.Equal(t, date_str, meal_date_str)
	assert.Equal(t, meal.UserID, user.ID)

	result = db.DB.Find(&meal_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(meal_images), 2)

	meal_image := meal_images[1]

	image_id := strconv.Itoa(meal_image.ID)
	filename := "meal/" + image_id + "/" + "demo1.png"
	file_url := "https://" + awsS3Bucket + ".s3." + awsRegion + ".amazonaws.com/" + filename
	assert.Equal(t, meal.ID, meal_image.MealID)
	assert.Equal(t, filename, meal_image.Filename)
	assert.Equal(t, file_url, meal_image.Fileurl)

	//

	//s3にあるか確認
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey, awsSecretKey, "",
		),
	}))

	download_file, _ := os.Create("downloader.png")
	defer download_file.Close()

	downloader := s3manager.NewDownloader(newSession)
	numBytes, err := downloader.Download(download_file,
		&s3.GetObjectInput{
			Bucket: aws.String(awsS3Bucket),
			Key:    aws.String(filename),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(numBytes)
	assert.Equal(t, true, numBytes > 0)

	os.Remove("downloader.png") //削除しておく

	//

	// json確認

	responseBody, _ := io.ReadAll(response2.Body)
	type ResponseData struct {
		Image model.MealImage
	}

	resp := ResponseData{}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, meal_image.ID, resp.Image.ID)
	assert.Equal(t, file_url, resp.Image.Fileurl)
}

func TestMealImageDelete_Success(t *testing.T) {
	//mealなしの場合はmealを作成してからmeal_imageを作成する
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsS3Bucket := os.Getenv("AWS_S3_TEST_BUCKET")

	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken

	//1度 meal,image作成
	access_token_list := []string{access_token}
	date_str_list := []string{"2022-07-24T00:00:00Z"}
	dummy_meal_create(access_token_list, date_str_list)

	menu_images := []model.MealImage{}
	result := db.DB.Find(&menu_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menu_images), 4)

	meal_image := menu_images[0]
	meal_image_id := meal_image.ID

	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/v2/meal/delete/meal_image/%v/", meal_image_id), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", access_token))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)

	result = db.DB.Find(&menu_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menu_images), 3)
	for _, m := range menu_images {
		assert.NotEqual(t, meal_image_id, m.ID)
	}

	//s3確認
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey, awsSecretKey, "",
		),
	}))

	s3.New(newSession)
	download_file, _ := os.Create("downloader.png")
	defer download_file.Close()

	downloader := s3manager.NewDownloader(newSession)
	_, err := downloader.Download(download_file,
		&s3.GetObjectInput{
			Bucket: aws.String(awsS3Bucket),
			Key:    aws.String(meal_image.Filename),
		},
	)

	assert.NotEqual(t, nil, err)

	os.Remove("downloader.png") //削除しておく
	//

	responseBody, _ := io.ReadAll(response.Body)
	type ResponseData struct {
		Meal_image model.MealImage
	}

	resp := ResponseData{}

	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, meal_image_id, resp.Meal_image.ID)

}

func TestMealImageDelete_FailDifferentUser(t *testing.T) {

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsS3Bucket := os.Getenv("AWS_S3_TEST_BUCKET")

	dbsetup()
	user := dummy_user_create()
	token := dummy_jwt_token_create(user.ID)
	access_token := token.AccessToken
	another_user := dummy_user_create2()
	another_token := dummy_jwt_token_create(another_user.ID)
	another_access_token := another_token.AccessToken

	//1度 meal,image作成
	access_token_list := []string{access_token}
	date_str_list := []string{"2022-07-24T00:00:00Z"}
	dummy_meal_create(access_token_list, date_str_list)

	menu_images := []model.MealImage{}
	result := db.DB.Find(&menu_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menu_images), 4)

	meal_image := menu_images[0]
	meal_image_id := meal_image.ID

	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/v2/meal/delete/meal_image/%v/", meal_image_id), nil)
	request.Header.Add("Authorization", fmt.Sprintf("JWT %v", another_access_token))
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code)

	result = db.DB.Find(&menu_images)
	if result.Error != nil {
		log.Fatalln(result.Error)
		return
	}
	assert.Equal(t, len(menu_images), 4)

	//s3にあるか確認
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey, awsSecretKey, "",
		),
	}))

	download_file, _ := os.Create("downloader.png")
	defer download_file.Close()

	downloader := s3manager.NewDownloader(newSession)
	numBytes, err := downloader.Download(download_file,
		&s3.GetObjectInput{
			Bucket: aws.String(awsS3Bucket),
			Key:    aws.String(meal_image.Filename),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, true, numBytes > 0)

	os.Remove("downloader.png") //削除しておく

	//

}
