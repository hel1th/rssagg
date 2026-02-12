FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/api/v1

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /app/server .

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["sh", "-c", "\
    until goose -dir ./migrations/schema postgres \"$DB_URL\" up; do \
    echo 'waiting for postgres...'; \
    sleep 2; \
    done && \
    ./server \
    "]
