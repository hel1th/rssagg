.PHONY: build run test clean docker-build docker-up docker-down migrate-up migrate-down sqlc

DB_URL ?= postgres://postgres:postgres@localhost:5432/rssagg?sslmode=disable

build:
	go build -o bin/api ./cmd/api/v1

run:
	go run ./cmd/api/v1

test:
	go test ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t rssagg .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	cd migration/schema && goose postgres "$(DB_URL)" up

migrate-down:
	cd migration/schema && goose postgres "$(DB_URL)" down

sqlc:
	sqlc generate
