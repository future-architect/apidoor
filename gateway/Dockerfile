FROM golang:1.17.7 as build
WORKDIR /gateway
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]
ADD ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w -buildid=' -trimpath -o /bin/main cmd/localredisgateway/main.go

FROM gcr.io/distroless/base
WORKDIR /log
EXPOSE 3000
WORKDIR /
COPY --from=build /bin/main .
CMD ["/main"]
