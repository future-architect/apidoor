module managementapi

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.3
	github.com/lib/pq v1.10.2
	local.packages/managementapi v0.0.0
)

replace local.packages/managementapi => ./
