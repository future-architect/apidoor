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

```bash
# Clone me
git clone https://github.com/future-architect/apidoor.git
cd apidoor

# Build all services
docker compose build \
  --build-arg http_proxy=${YOUR_PROXY} \
  --build-arg https_proxy=${YOUR_PROXY} \
  --build-arg proxy=${YOUR_PROXY} \
  --build-arg https-proxy=${YOUR_PROXY}

# Launch apidoor services
docker compose up -d

# Set your first API routing through management-api
curl -X POST -H "Content-Type: application/json" \
-d '{"api_key": "key", "path": "test", "forward_url": "http://test-server:3333/welcome"}' localhost:3001/mgmt/api

# Check apidoor works
curl -H "Content-Type: application/json" -H "Authorization:key" localhost:3000/test
# welcome to apidoor!

# Check log file is provided
cat ./log/log.csv

# You can also access Management Console
localhost:8080
```

## Architecture

TODO

# License
Apache 2
