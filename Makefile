run:
	go run ./cmd/api

build:
	go build -o main ./cmd/api

docker-up:
	docker compose up -d

docker-down:
	docker compose down