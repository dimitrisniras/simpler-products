install-deps:
	go get github.com/gin-gonic/gin
	go get github.com/joho/godotenv
	go get github.com/go-sql-driver/mysql
	go get github.com/sirupsen/logrus
	go get github.com/go-playground/validator/v10

build:
	go build -o bin/main main.go

run:
	go run main.go