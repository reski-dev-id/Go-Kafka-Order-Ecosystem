run:
	cd service-order && go run ./cmd/api

build:
	cd service-order && go build -o main ./cmd/api

swagger:
	cd service-order && swag init -g cmd/api/main.go

docker-up:
	cd service-order && docker compose up -d

docker-down:
	cd service-order && docker compose down