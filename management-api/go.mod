module managementapi

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/go-chi/chi/v5 v5.0.3
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/swaggo/swag v1.7.0
	local.packages/managementapi v0.0.0
)

replace local.packages/managementapi => ./
