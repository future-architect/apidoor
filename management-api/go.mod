module managementapi

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-redis/redis/v8 v8.11.3 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/swaggo/swag v1.7.0
	go.opentelemetry.io/otel v0.14.0 // indirect
	local.packages/managementapi v0.0.0
)

replace local.packages/managementapi => ./
