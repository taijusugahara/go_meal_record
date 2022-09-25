package main

import (
	"github.com/joho/godotenv"
	"go-meal-record/app/db" //module名/ディレクトリ initだけの場合は_使用(下の理由によりinitはやめた)
	"go-meal-record/app/router"
	"log"
	"os"
	"time"
)

func main() {
	//circleciを利用する=githubを利用するということは.envは使えないので、開発環境だけ.env使って、本番環境(ecs)はタスク定義に環境変数を書く
	if os.Getenv("GO_ENVIRONMENT") == "development" { //productionでは.envないからエラーなる
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln(err)
		}
	}
	db.SettingDb() //元々dbのsettingはinitにしていたが、initだとtest環境でenvを上書きする(test環境だと伝える)前に実行されてしまい、開発環境と同じdbが使用されてしまうのでinitを使うのはやめた。
	engine := router.Router()
	log.SetFlags(log.Lshortfile) //デフォルトだと何行目かの情報でないので行情報出るように。
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	engine.Run(":3000")
}
