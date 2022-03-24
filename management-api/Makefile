#################
# Unit Tests    #
#################
test: test-redis test-dynamo

test-redis:
	export TEST_MODE="test"; \
	export DATABASE_DRIVER="postgres"; \
	export DATABASE_HOST="127.0.0.1"; \
	export DATABASE_PORT="5432"; \
	export DATABASE_USER="root"; \
	export DATABASE_PASSWORD="password"; \
	export DATABASE_NAME="root"; \
	export DATABASE_SSLMODE="disable"; \
	export API_DB_TYPE="REDIS"; \
	export REDIS_HOST="localhost"; \
	export REDIS_PORT="6379"; \
	go test -race -v ./

test-dynamo:
	export TEST_MODE="test"; \
	export DATABASE_DRIVER="postgres"; \
	export DATABASE_HOST="127.0.0.1"; \
	export DATABASE_PORT="5432"; \
	export DATABASE_USER="root"; \
	export DATABASE_PASSWORD="password"; \
	export DATABASE_NAME="root"; \
	export DATABASE_SSLMODE="disable"; \
	export API_DB_TYPE="DYNAMO"; \
	export DYNAMO_TABLE_API_ROUTING="api_routing"; \
	export DYNAMO_TABLE_ACCESS_TOKEN="access_token"; \
	export DYNAMO_TABLE_SWAGGER="swagger"; \
	export DYNAMO_ENDPOINT="http://localhost:4566"; \
	go test -race -v ./...

#################
# swagger       #
#################
generate-swagger:
	swag init -g cmd/management-api/main.go

