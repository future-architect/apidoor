FROM golang:1.16.4

WORKDIR /go/src/cmd

COPY cmd /go/src/cmd/

COPY *.go go.* /go/src/

WORKDIR /go/src

RUN ["go", "mod", "download"]

WORKDIR /go/src/cmd/main

CMD ["go", "run", "main.go"]