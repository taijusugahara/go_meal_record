package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-meal-record/app/db"
	"go-meal-record/app/router"
	"log"
	"os"
	"testing"
	"time"
)

var engine *gin.Engine

//このTestMainは.._test.goに入れないとと動かないみたい
func TestMain(m *testing.M) { //テスト前と後の共通処理 m.Run()の前後で
	start_test()
	m.Run()
	aws_s3_file_all_delete()
}

func start_test() {
	is_circleci_test := os.Getenv("IS_CIRCLECI_TEST")
	if is_circleci_test != "true" { //circleciでは.env使わずcircleciの環境変数使う。
		os.Setenv("GO_ENVIRONMENT", "test")
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("this is circleci test")
	}

	engine = router.Router()
	log.SetFlags(log.Lshortfile) //デフォルトだと何行目かの情報でないので行情報出るように。
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
}

func aws_s3_file_all_delete() {
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsS3Bucket := os.Getenv("AWS_S3_TEST_BUCKET")

	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey, awsSecretKey, "",
		),
	}))

	svc := s3.New(newSession)

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(awsS3Bucket),
	})

	err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter)
	if err != nil {
		log.Fatalln(err)
	}
}

func dbsetup() {
	//goはtest終わってもdb保持されてしまうため、test1回１回始まる度db初期化
	db.SettingDb()
}
