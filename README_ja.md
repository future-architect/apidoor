![logo](docs/apidoor_logo.png)

# apidoor

apidoor は API の商材管理や利用状況確認を便利にする OSS です。

## What is apidoor for

TODO

## Features

## Prerequisites
- go v1.16^
- redis-server v6.2^
- docker v20.10^
- docker-compose v1.29^
- npm v6.14^

## Getting Started

### リポジトリのクローン
```
git clone https://gitlab.com/osaki-lab/apidoor.git
cd apidoor
```

### redis-server の起動
API エンドポイントや利用回数の管理のため、redis-server を起動します。
```
sudo service redis-server start
```

### gateway の起動
管理している API へのリクエストをプロキシするゲートウェイを起動します。環境変数などの細かい設定に関しては[こちら](https://gitlab.com/osaki-lab/apidoor/-/tree/master/gateway)をご覧ください。
```
cd gateway/cmd/gateway
go run main.go
```

### management-api の起動
管理している API に関する情報やその利用状況を提供する API を起動します。環境変数などの細かい設定に関しては[こちら](https://gitlab.com/osaki-lab/apidoor/-/tree/master/management-api)をご覧ください。
```
cd management-api
# proxy環境下での実行時のみ
docker-compose build \
  --build-arg HTTP_PROXY=${YOUR_PROXY} \
  --build-arg HTTPS_PROXY=${YOUR_PROXY} \
  --build-arg http_proxy=${YOUR_PROXY} \
  --build-arg ${YOUR_PROXY}

docker-compose up -d
```

### management-front の起動
GUI 上で API の管理等を行う Vue.js アプリケーションを起動します。環境変数などの細かい設定に関しては[こちら](https://gitlab.com/osaki-lab/apidoor/-/tree/master/management-front)をご覧ください。
```
npm install
npm run serve
```
コマンドの実行後、[localhost:8081](localhost:8081) にアクセスしてください。

## Architecture

TODO

# License
Apache 2

