package controller

import (
	"go-meal-record/app/db"
	"go-meal-record/app/model"
	"go-meal-record/app/utils/common"
	"go-meal-record/app/utils/validate"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func MealIndexByDay(c *gin.Context) {
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID

	date_str := c.Param("date")
	//2006-01-.....これじゃないといけないみたい
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	//dateはmodelでdatatypes.Date使用しているので例えば2022:07:20T10:40:30Zは2022:07:20T:00:00:00Zとなる。
	//そのため時間は00:00:00,0固定でいい。
	meals := []model.Meal{}
	one_day_meal_map := map[string]interface{}{
		"morning": model.NewMeal(), "lunch": model.NewMeal(), "dinner": model.NewMeal(), "other": model.NewMeal(),
	} //注意 : mapは強制的にkey abc順になる、morningから開始したいがslice rangeで順番変えても無理だった。今回の場合はdinnerが先頭に来る。morningからにしたかったらjs側でsortするか

	result := db.DB.Model(model.Meal{}).Where("date = ? AND user_id = ?", date_start, user_id).Preload("Menus").Preload("MealImages").Joins("User").Order("date").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	for _, meal := range meals {
		meal_type := meal.MealType
		one_day_meal_map[meal_type] = meal
	}

	c.JSONP(http.StatusOK, gin.H{
		"data": one_day_meal_map,
	})
}

func MealIndexByWeek(c *gin.Context) {
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID

	meals := []model.Meal{}
	date_str := c.Param("date")
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	weekday := date.Weekday()
	weekday_int := int(weekday)
	first_week_date := date_start.AddDate(0, 0, -weekday_int)
	end_week_date := date_start.AddDate(0, 0, 6-weekday_int)
	formatted_first_week_date := first_week_date.Format("2006-01-02T15:04:05Z")
	formatted_end_week_date := end_week_date.Format("2006-01-02T15:04:05Z")
	result := db.DB.Where("meals.date BETWEEN ? AND ? AND user_id= ?", first_week_date, end_week_date, user_id).Preload("Menus").Preload("MealImages").Joins("User").Order("date").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	one_week_meal_map := map[string]map[string]interface{}{}
	for i := 0; i < 7; i++ { //1週間分の日時をkeyにして初期値を入れる
		weekday_date := first_week_date.AddDate(0, 0, i)
		weekday_date_str := weekday_date.Format("2006-01-02T15:04:05Z")
		one_day_meal_map := map[string]interface{}{
			"morning": model.NewMeal(), "lunch": model.NewMeal(), "dinner": model.NewMeal(), "other": model.NewMeal(),
		}
		one_week_meal_map[weekday_date_str] = one_day_meal_map
	}

	for _, meal := range meals {
		date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
		date_str := date.Format("2006-01-02T15:04:05Z")
		meal_type := meal.MealType
		one_week_meal_map[date_str][meal_type] = meal
	}
	c.JSONP(http.StatusOK, gin.H{
		"data":            one_week_meal_map,
		"first_week_date": formatted_first_week_date,
		"end_week_date":   formatted_end_week_date,
	})
}

func MealIndexByMonth(c *gin.Context) {
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID

	meals := []model.Meal{}
	date_str := c.Param("date")
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	first_month_date := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	end_month_date := first_month_date.AddDate(0, 1, -1)
	formatted_first_month_date := first_month_date.Format("2006-01-02T15:04:05Z")
	result := db.DB.Where("meals.date BETWEEN ? AND ? AND user_id= ?", first_month_date, end_month_date, user_id).Preload("Menus").Preload("MealImages").Joins("User").Order("id").Find(&meals)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	one_month_meal_map := map[string]map[string]interface{}{}
	this_month := first_month_date.Month()
	for i := 0; ; i++ { //1ヶ月分の日時をkeyにして初期値を入れる
		month_date := first_month_date.AddDate(0, 0, i)
		if month_date.Month() != this_month {
			break
		}
		month_str := month_date.Format("2006-01-02T15:04:05Z")
		one_day_meal_map := map[string]interface{}{
			"morning": model.NewMeal(), "lunch": model.NewMeal(), "dinner": model.NewMeal(), "other": model.NewMeal(),
		}
		one_month_meal_map[month_str] = one_day_meal_map
	}
	for _, meal := range meals {
		date := time.Time(meal.Date) //1回time.Time型にしてあげる。UTCだった。これをjstにすると+9時間されてしまうのでここではしない。
		date_str := date.Format("2006-01-02T15:04:05Z")
		meal_type := meal.MealType
		one_month_meal_map[date_str][meal_type] = meal
	}
	c.JSONP(http.StatusOK, gin.H{
		"data":             one_month_meal_map,
		"first_month_date": formatted_first_month_date,
	})
}

func MealGetOrCreate(c *gin.Context, meal_type string, date_str string) (model.Meal, error) {
	validate := validate.Validate()
	var err error
	meal := model.Meal{}
	already_created_meal := model.Meal{}
	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return meal, err
	}
	user_id := user.ID
	date, _ := time.ParseInLocation("2006-01-02T15:04:05Z", date_str, time.Local)
	date_start := datatypes.Date(date)
	// ErrRecordNotFound エラーを避けたい場合は、db.Limit(1).Find(&user)のように、Find を 使用することができます。
	db.DB.Where("meal_type = ? AND user_id = ? AND date = ?", meal_type, user_id, date_start).Joins("User").Limit(1).Find(&already_created_meal)
	if already_created_meal.ID != 0 {
		meal = already_created_meal
	} else { //新しく作成する
		meal.MealType = meal_type
		meal.Date = date_start
		meal.User = user
		meal.UserID = user_id

		err = validate.Struct(meal)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Validation Error")
			return meal, err
		}

		result := db.DB.Create(&meal)
		if result.Error != nil {
			log.Println(result.Error)
			c.String(http.StatusInternalServerError, "Server Error")
			return meal, result.Error
		}
	}

	return meal, nil
}

func MenuCreate(c *gin.Context) {
	validate := validate.Validate()
	var err error
	menu := model.Menu{}
	//js FormDataから渡してる。JsonBindしないのは渡ってきた3つが1つのstructに入るわけではないから。
	meal_type := c.PostForm("meal_type")
	date_str := c.PostForm("date")
	menu_name := c.PostForm("menu")

	meal, err := MealGetOrCreate(c, meal_type, date_str)
	if err != nil {
		return
	}
	//同じ名前とmeal_idを持つmenuは作成させない。
	already_created_menu := model.Menu{}
	// ErrRecordNotFound エラーを避けたい場合は、db.Limit(1).Find(&user)のように、Find を 使用することができます。
	db.DB.Where("name = ? AND meal_id = ?", menu_name, meal.ID).Limit(1).Find(&already_created_menu)
	if already_created_menu.ID != 0 { //存在したら
		c.String(http.StatusInternalServerError, "既に同じ名前とmeal_idを持つmenuが存在するため作成できません。")
		return
	}

	menu.Name = menu_name
	menu.Meal = meal
	menu.MealID = meal.ID
	err = validate.Struct(menu)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Validation Error")
		return
	}
	result := db.DB.Create(&menu)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"menu": menu,
	})
}

func MenuDelete(c *gin.Context) {
	id := c.Param("id") //dbにidを渡す際、stringでもintでもどっちもでもいいみたい。
	menu := model.Menu{}
	result := db.DB.Joins("Meal").First(&menu, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID
	if menu.Meal.UserID != user_id {
		c.String(http.StatusBadRequest, "User not correct user")
		return
	}
	result = db.DB.Delete(&menu)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"menu": menu,
	})

}

func MealImageCreate(c *gin.Context) {
	var err error
	validate := validate.Validate()
	image := model.MealImage{}
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	go_environment := os.Getenv("GO_ENVIRONMENT")
	awsS3Bucket := os.Getenv("AWS_S3_LOCAL_BUCKET")
	if go_environment == "production" {
		awsS3Bucket = os.Getenv("AWS_S3_PRODUCTION_BUCKET")
	} else if go_environment == "test" {
		awsS3Bucket = os.Getenv("AWS_S3_TEST_BUCKET")
	}
	//fileがある場合はjsonBindじゃなくてc.PostForm使っていけた。
	//理由としてはデータ格納するのにreact側でFormData使ってるからだと思う
	meal_type := c.PostForm("meal_type")
	date_str := c.PostForm("date")
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}
	//Menuがあるか確認。なければ作成
	meal, err := MealGetOrCreate(c, meal_type, date_str)

	filename := file.Filename
	filename_split_dot := strings.Split(filename, ".")
	extention := filename_split_dot[len(filename_split_dot)-1]
	valid_extentions := []string{"jpeg", "jpg", "JPEG", "png", "PNG"}
	if common.Contains(valid_extentions, extention) {
		// image.Filename = filename
		image.Meal = meal
		image.MealID = meal.ID
		err = validate.Struct(image)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Validation Error")
			return
		}
		//一旦pkを獲得するためにcreate
		result := db.DB.Create(&image)
		if result.Error != nil {
			log.Println(result.Error)
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}
		image_id := strconv.Itoa(image.ID)
		filename := "meal/" + image_id + "/" + filename
		// sessionを作成します(aws 接続)
		newSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(
				awsAccessKey, awsSecretKey, "",
			),
		}))
		upload_file, _ := file.Open()
		defer upload_file.Close()

		// Uploaderを作成し、ファイルをアップロード
		uploader := s3manager.NewUploader(newSession)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(awsS3Bucket),
			Key:    aws.String(filename),
			Body:   upload_file,
		})

		if err != nil {
			db.DB.Delete(&image)
			log.Print(err)
			c.String(http.StatusInternalServerError, "s3 upload fail")
			return
		}
		// file object更新
		image.Filename = filename
		image.Fileurl = "https://" + awsS3Bucket + ".s3." + awsRegion + ".amazonaws.com/" + filename
		db.DB.Save(&image)

	} else {
		c.String(http.StatusInternalServerError, "File extention not correct")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"image": image,
	})
}

func MealImageDelete(c *gin.Context) {
	id := c.Param("id") //dbにidを渡す際、stringでもintでもどっちもでもいいみたい。
	meal_image := model.MealImage{}
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	go_environment := os.Getenv("GO_ENVIRONMENT")
	awsS3Bucket := os.Getenv("AWS_S3_LOCAL_BUCKET")
	if go_environment == "production" {
		awsS3Bucket = os.Getenv("AWS_S3_PRODUCTION_BUCKET")
	} else if go_environment == "test" {
		awsS3Bucket = os.Getenv("AWS_S3_TEST_BUCKET")
	}

	result := db.DB.Joins("Meal").First(&meal_image, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	user, err := GetCurrentUser(c)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "Not correct user")
		return
	}
	user_id := user.ID
	if meal_image.Meal.UserID != user_id {
		c.String(http.StatusBadRequest, "User not correct user")
		return
	}
	filename := meal_image.Filename

	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey, awsSecretKey, "",
		),
	}))

	svc := s3.New(newSession)
	// s3はフォルダ内のファイル全て削除されたらフォルダも自動で削除されるらしい。なのでファイル削除だけでok
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(awsS3Bucket), Key: aws.String(filename),
	})
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsS3Bucket),
		Key:    aws.String(filename),
	}) // errorが帰ってこなければ成功

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	result = db.DB.Delete(&meal_image)
	if result.Error != nil {
		log.Println(result.Error)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"meal_image": meal_image,
	})

}
