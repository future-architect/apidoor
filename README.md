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
* [ ] Auto publish of an API access token
* [ ] Management of products
* [ ] Check the usage situation of APIs

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

