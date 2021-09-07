module gateway

go 1.16

require (
	github.com/aws/aws-sdk-go v1.40.37 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/guregu/dynamo v1.11.0 // indirect
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f // indirect
	local.packages/gateway v0.0.0-00010101000000-000000000000
)

replace local.packages/gateway => ./
