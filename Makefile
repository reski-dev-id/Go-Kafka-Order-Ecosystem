run-order:
	cd service-order && go run ./cmd/api

run-relay:
	cd relay-worker && go run ./cmd/worker/main.go

run-payment:
	cd payment-service && ./mvnw spring-boot:run

run-all:
	@echo "starting order-service..."
	cd service-order && nohup go run ./cmd/api > ../order-service.log 2>&1 &

	@echo "starting relay-worker..."
	cd relay-worker && nohup go run ./cmd/worker/main.go > ../relay-worker.log 2>&1 &

	@echo "starting payment-service..."
	cd payment-service && nohup ./mvnw spring-boot:run > ../payment-service.log 2>&1 &

	@echo "all services started"

logs:
	tail -f order-service.log relay-worker.log payment-service.log

stop-all:
	-pkill -9 -f "/home/reski/.cache/go-build"
	-pkill -9 -f "PaymentServiceApplication"

build:
	cd service-order && go build -o main ./cmd/api

swagger:
	cd service-order && swag init -g cmd/api/main.go

docker-up:
	docker compose up -d

docker-down:
	docker compose down