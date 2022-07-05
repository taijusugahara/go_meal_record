# image取得(容量が小さいalpineを選択)
FROM golang:1.17.6-alpine

# ホストのファイルをコンテナの作業ディレクトリにコピー
COPY . /go/src

# ワーキングディレクトリの設定
#上で設定したpathから考える
WORKDIR /go/src/app