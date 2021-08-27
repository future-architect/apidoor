![logo](docs/apidoor_logo.png)

# apidoor

apidoor は API の商材管理や利用状況確認を便利にする OSS です。

## What is apidoor for

TODO

## Features

## Prerequisites
- 
- Go v1.16^
- redis-server v6.2^
- docker v20.10^
- docker-compose v1.29^
- npm v6.14^

## Getting Started

```
# Clone me
git clone https://gitlab.com/osaki-lab/apidoor.git
cd apidoor

# Build all services
docker-compose build \
  --build-arg http_proxy=${YOUR_PROXY} \
  --build-arg https_proxy=${YOUR_PROXY}

# Launch apidoor services
docker compose up -d

# Set your first API routing
docker exec -it redis-server sh
> redis-cli
127.0.0.1:6379> hset key test test-server:3333/welcome
127.0.0.1:6379> exit
> exit

# Check apidoor works
curl -H "Content-Type: application/json" -H "Authorization: key" localhost:3000/test
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
