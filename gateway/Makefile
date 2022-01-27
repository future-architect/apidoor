#################
# Unit Tests    #
#################
unittest:
	export READTIMEOUT=5; \
	export READHEADERTIMEOUT=5; \
	export WRITETIMEOUT=20; \
	export IDLETIMEOUT=5; \
	export MAXHEADERBYTES="1<<20"; \
	export DYNAMO_TABLE_ACCESS_LOG="access_log"; \
	export DYNAMO_ACCESS_LOG_ENDPOINT="http://localhost:4566"; \
	export LOG_PATH="./testdata/log/log.csv"; \
	export API_DB_TYPE="REDIS"; \
	go test -race -v ./...

run:
	go run cmd/localredisgateway/main.go

test:
	go test ./...
