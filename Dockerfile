
FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache make

RUN make install-deps

RUN make build

EXPOSE 8080

CMD ["./bin/main"]