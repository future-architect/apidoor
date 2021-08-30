![logo](docs/apidoor_logo.png)

# apidoor

apidoor is an OSS product that accelerates the construction of the API market.

[日本語](README_ja.md)

**🚧This project is Work in Progress🚧**

## What is apidoor for
You can use apidoor when

* you want to open APIs for many users and have trouble to publish a list of API and an access token.

## Features

* [x] Routing and access management of WebAPI
* [ ] Auto publishment of an API access token
* [ ] Management of products
* [ ] Check the usage situation of APIs

## Getting Started

Prerequisites:

- docker v20.10^
- docker-compose v1.29^

Flow：

```
# Clone me
git clone https://gitlab.com/osaki-lab/apidoor.git
cd apidoor

# fix docker-compose
## comment out L24 or L26 depending on your OS
vi docker-compose.yml

# Build all services
docker-compose build \
  --build-arg http_proxy=${YOUR_PROXY} \
  --build-arg https_proxy=${YOUR_PROXY} \
  --build-arg proxy=${YOUR_PROXY} \
  --build-arg https-proxy=${YOUR_PROXY}

# Launch apidoor services
docker-compose up -d

# Set your first API routing
docker exec -it redis-server sh
> redis-cli
127.0.0.1:6379> hset key test test-server:3333/welcome
127.0.0.1:6379> exit
> exit

# Check apidoor works
curl -H "Content-Type: application/json" -H "Authorization:key" localhost:3000/test
# welcome to apidoor!

# Check log file is provided
## for mac
sudo cat log/log.csv
## for WSL
sudo cat /mnt/c/log/log.csv

# You can also access Management Console
localhost:8080
```

## Architecture

TODO

# License
Apache 2

