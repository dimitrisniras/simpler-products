# Project variables
PROJECT_NAME = simpler-products-app
GO_VERSION = 1.23

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

# Phony targets
.PHONY: all install-deps build run test test-cover clean clean-deps lint fmt docker-build docker-run docker-stop

all: clean-all install-deps build run

install-deps:
	go get github.com/gin-gonic/gin
	go get github.com/joho/godotenv
	go get github.com/go-sql-driver/mysql
	go get github.com/sirupsen/logrus
	go get github.com/google/uuid
	go get github.com/golang-jwt/jwt/v4
	go get github.com/go-playground/validator/v10
	go get github.com/dgrijalva/jwt-go
	go get github.com/DATA-DOG/go-sqlmock
	go get github.com/stretchr/testify/assert
	go get github.com/dvwright/xss-mw
	go get github.com/golangci/golangci-lint

build:
	go build -o bin/main main.go

run:
	go run main.go

test:
	go test ./tests -v

# Cleaning
clean:
	rm -rf bin/*

clean-deps:
	go clean -modcache

clean-all: clean clean-deps

# Linting and formatting
lint: 
	${GOLANGCI_LINT} run --no-config --disable-all --enable gocritic,gofumpt

lint-fix: 
	${GOLANGCI_LINT} run --no-config --disable-all --enable gocritic,gofumpt --fix

fmt:
	go fmt ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Docker
docker-build:
	docker build -t ${PROJECT_NAME} .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down
