FROM golang:1.18
COPY . /f1-example
WORKDIR /f1-example
CMD go run ./cmd/service/main.go
