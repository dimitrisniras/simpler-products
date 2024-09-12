install-deps:
	go get github.com/gin-gonic/gin
	go get github.com/joho/godotenv
	go get github.com/go-sql-driver/mysql

build:
	go build -o bin/main main.go

run:
	go run main.go