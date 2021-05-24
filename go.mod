module apidoor

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.3
	local.packages/apidoor v0.0.0-00010101000000-000000000000
)

replace local.packages/apidoor => ./
