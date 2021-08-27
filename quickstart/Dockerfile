FROM golang:1.16.4 as build
WORKDIR /quickstart
ADD cmd/*.go cmd/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w -buildid=' -trimpath -o /bin/main cmd/main.go

FROM gcr.io/distroless/base
EXPOSE 3333
COPY --from=build /bin/main .
CMD ["/main"]
