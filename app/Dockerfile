# image取得(容量が小さいalpineを選択)
FROM golang:1.17.6-alpine

# ホストのファイルをコンテナの作業ディレクトリにコピー
# 下の.はdocker-composeでcontextを.としているためdocker-composeのある階層になってる
COPY ./app /go/src/go_meal_record/app

# ワーキングディレクトリの設定
WORKDIR /go/src/go_meal_record/app
