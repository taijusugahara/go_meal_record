# image取得(容量が小さいalpineを選択)
FROM golang:1.17.6-alpine

# ホストのファイルをコンテナの作業ディレクトリにコピー
COPY . /go/src/app

# ワーキングディレクトリの設定
#copyとwordirはこの形が一番かな？appのところをgo_meal_recordとかにしたらおかしくなる
WORKDIR /go/src/app