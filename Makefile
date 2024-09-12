install-deps:
	go get github.com/gin-gonic/gin
	go get github.com/joho/godotenv

build:
	go build -o bin/main main.go

run:
	go run main.go