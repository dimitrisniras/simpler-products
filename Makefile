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

build:
	go build -o bin/main main.go

run:
	go run main.go

test:
	go test ./tests -v