module redislogger

go 1.16

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-redis/redis/v8 v8.11.3
	go.opentelemetry.io/otel v0.20.0 // indirect
)

replace local.packages/redislogger => ./
