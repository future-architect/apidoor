#################
# Unit Tests    #
#################
redis_test:
	echo "===== Redis test ====="
	export READTIMEOUT=5; \
	export READHEADERTIMEOUT=5; \
	export WRITETIMEOUT=20; \
	export IDLETIMEOUT=5; \
	export MAXHEADERBYTES="1<<20"; \
	export DB_TYPE="REDIS"; \
	export REDIS_HOST="localhost"; \
	export REDIS_PORT="6379"; \
	export LOG_PATH="./testdata/log/log.csv"; \
	go test -race -v ./

dynamodb_test:
	echo "===== Dynamo DB test ====="
	export READTIMEOUT=5; \
	export READTIMEOUT=5; \
	export READHEADERTIMEOUT=5; \
	export WRITETIMEOUT=20; \
	export IDLETIMEOUT=5; \
	export MAXHEADERBYTES="1<<20"; \
	export DB_TYPE="DYNAMO"; \
	export DYNAMO_TABLE_API_FORWARDING="local_api_forwarding"; \
	export DYNAMO_ENDPOINT="http://localhost:4566"; \
	export AWS_REGION=ap-northeast-1; \
	export AWS_DEFAULT_OUTPUT=text; \
	export AWS_ACCESS_KEY_ID=xxx; \
	export AWS_SECRET_ACCESS_KEY=xxx; \
	export LOG_PATH="./testdata/log/log.csv"; \
	aws dynamodb --endpoint-url http://localhost:4566 create-table --cli-input-json file://testdata/model/local_api_forwarding.json; \
	go test -race -v ./; \
	aws dynamodb --endpoint-url http://localhost:4566 delete-table --table local_api_forwarding

test: redis_test dynamodb_test
