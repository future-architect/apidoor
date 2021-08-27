![logo](docs/apidoor_logo.png)

# apidoor

apidoor は WebAPI の商材管理や利用状況確認を便利にする OSS です。

## What is apidoor for

* WebAPIを様々な接続先に公開するときに、APIの一覧やアクセストークンの払い出しを行いたいユースケースの場合

## Features

* [x] WebAPIのルーティングやアクセス制限
* [ ] APIトークンの自動払い出し
* [ ] 商材管理
* [ ] 利用状況の確認

## Getting Started

Prerequisites:

- docker v20.10^
- docker-compose v1.29^

Flow：

```
# Clone me
git clone https://gitlab.com/osaki-lab/apidoor.git
cd apidoor

# Build all services
docker-compose build \
  --build-arg http_proxy=${YOUR_PROXY} \
  --build-arg https_proxy=${YOUR_PROXY} \
  --build-arg proxy=${YOUR_PROXY} \
  --build-arg https-proxy=${YOUR_PROXY}

# Launch apidoor services
docker compose up -d

# Set your first API routing
docker exec -it redis-server sh
> redis-cli
127.0.0.1:6379> hset key test test-server:3333/welcome
127.0.0.1:6379> exit
> exit
> exit

# Check apidoor works
curl -H "Content-Type: application/json" -H "Authorization:key" localhost:3000/test
# welcome to apidoor!

# You can also access Management Console
localhost:8080
```

## Architecture

TODO

# License
Apache 2
