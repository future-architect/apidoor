FROM golang:1.17.7 as build
WORKDIR /management-api
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]
ADD ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w -buildid=' -trimpath -o /bin/main cmd/management-api/main.go

FROM gcr.io/distroless/base-debian10:latest
EXPOSE 3001
COPY --from=build /bin/main .
CMD ["/main"]
