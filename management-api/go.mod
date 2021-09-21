module managementapi

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/google/go-cmp v0.5.6
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/swaggo/swag v1.7.1
	gopkg.in/go-playground/validator.v8 v8.18.2
	local.packages/managementapi v0.0.0
)

replace local.packages/managementapi => ./
