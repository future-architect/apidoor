#################
# Unit Tests    #
#################
test:
	export READTIMEOUT=5; \
	export READHEADERTIMEOUT=5; \
	export WRITETIMEOUT=20; \
	export IDLETIMEOUT=5; \
	export MAXHEADERBYTES="1<<20"; \
	export LOG_PATH="./testdata/log/log.csv"; \
	go test -race -v ./

run:
	go run cmd/localredisgateway/main.go
