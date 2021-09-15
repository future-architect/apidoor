#################
# Unit Tests    #
#################
test:
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


#################
# Unit Tests    #
#################
generate-swagger:
	swag init -g cmd/management-api/main.go

