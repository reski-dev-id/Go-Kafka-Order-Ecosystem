# Go Kafka Order Ecosystem

Production-style distributed system implementing **Event-Driven Architecture** with the **Transactional Outbox Pattern** to guarantee consistency between PostgreSQL transactions and Kafka event publishing.

This project demonstrates:

- Microservice architecture
- Polyglot services
- Clean Architecture
- Feature-based structure
- Dependency Injection
- Kafka event streaming
- Distributed tracing
- Metrics & observability
- Dockerized infrastructure
- CI/CD with GitHub Actions
- DockerHub deployment

---

# Architecture

```text
Client
  ↓
Order Service (Go + Gin)
  ↓
PostgreSQL (orders + outbox)
  ↓
Relay Worker (Go)
  ↓
Kafka
  ↓
Payment Service (Spring Boot)
  ↓
Kafka
  ↓
Notification Service (FastAPI)
```

---

# Tech Stack

| Layer | Technology |
|---|---|
| Order Service | Go 1.22 + Gin |
| Relay Worker | Go 1.22 |
| Payment Service | Java Spring Boot |
| Notification Service | FastAPI |
| Database | PostgreSQL |
| SQL Library | sqlx |
| Message Broker | Apache Kafka |
| Metrics | Prometheus |
| Dashboard | Grafana |
| Logs | Loki |
| Tracing | OpenTelemetry + Tempo |
| Dependency Injection | Google Wire |
| Containerization | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Image Registry | DockerHub |

---

# Core Architecture Principles

## Microservice Architecture

Each service has a single responsibility:

| Service | Responsibility |
|---|---|
| order-service | Create orders + transactional outbox |
| relay-worker | Publish outbox events to Kafka |
| payment-service | Consume order events and process payment |
| notification-service | Send notifications from Kafka events |

---

## Clean Architecture

Every service follows Clean Architecture:

```text
┌─────────────────────────────────────────────┐
│               Delivery Layer                │
├─────────────────────────────────────────────┤
│               Usecase Layer                 │
├─────────────────────────────────────────────┤
│              Repository Layer               │
├─────────────────────────────────────────────┤
│               Domain Layer                  │
└─────────────────────────────────────────────┘
```

Rules:

- Domain layer must remain framework-independent
- Usecase layer must not depend on HTTP
- Repository interfaces belong to the feature layer
- Infrastructure only implements contracts
- All dependencies are injected

---

## Feature-Based Structure

```text
features/
└── order/
    ├── delivery/
    ├── usecase/
    ├── repository/
    ├── domain/
    └── dto/
```

Advantages:

- Better scalability
- Better ownership boundaries
- Easier maintenance
- Reduced coupling
- Cleaner modularity

---

# Repository Structure

```text
go-kafka-order-ecosystem/
│
├── contracts/
│   └── events/
│
├── shared/
│   ├── logger/
│   ├── telemetry/
│   ├── postgres/
│   ├── kafka/
│   ├── middleware/
│   ├── response/
│   └── errors/
│
├── order-service/
│   ├── cmd/
│   └── internal/
│
├── relay-worker/
│
├── payment-service/
│
├── notification-service/
│
├── infra/
│   ├── docker-compose.yml
│   ├── prometheus/
│   ├── grafana/
│   ├── loki/
│   ├── tempo/
│   └── postgres-init/
│
├── .github/
│   └── workflows/
│
└── README.md
```

---

# Dependency Injection

Dependency Injection uses:

- Google Wire

Benefits:

- Compile-time dependency graph
- Type-safe DI
- Better testing
- Cleaner initialization
- Explicit dependencies

---

# Database Schema

## Orders Table

```sql
CREATE TABLE orders (
  id UUID PRIMARY KEY,
  item_name VARCHAR(255),
  amount DECIMAL,
  status VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Outbox Table

```sql
CREATE TABLE outbox (
  id BIGSERIAL PRIMARY KEY,
  aggregate_id UUID NOT NULL,
  event_key VARCHAR(255),
  topic VARCHAR(100) NOT NULL,
  payload JSONB NOT NULL,

  processed BOOLEAN DEFAULT FALSE,

  retry_count INT DEFAULT 0,
  next_retry_at TIMESTAMP NULL,

  last_error TEXT,

  processed_at TIMESTAMP NULL,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Processed Events Table

Used for idempotent Kafka consumers.

```sql
CREATE TABLE processed_events (
  event_id VARCHAR(255) PRIMARY KEY,
  processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

# Kafka Topics

```text
order.created
payment.completed
payment.failed
notification.send
order.dlq
```

---

# Event Contract

All services communicate through a shared event schema.

```json
{
  "event_id": "uuid",
  "event_type": "order.created",
  "aggregate_id": "uuid",
  "occurred_at": "timestamp",
  "payload": {
    "order_id": "uuid",
    "amount": 100
  }
}
```

---

# Transactional Outbox Flow

## Step 1 — Order Creation

The order-service:

- inserts order data
- inserts outbox event
- commits in a single PostgreSQL transaction

---

## Step 2 — Relay Worker

The relay-worker:

- polls unprocessed outbox rows
- publishes events to Kafka
- marks rows as processed
- retries failed events

---

## Step 3 — Payment Service

The payment-service:

- consumes Kafka events
- validates idempotency
- processes payment
- publishes payment result events

---

## Step 4 — Notification Service

The notification-service:

- consumes payment events
- sends notifications

---

# Reliable Outbox Query

Relay worker uses:

```sql
SELECT *
FROM outbox
WHERE processed = false
ORDER BY created_at
LIMIT 10
FOR UPDATE SKIP LOCKED;
```

This prevents:

- duplicate processing
- worker race conditions
- row contention

---

# Observability Stack

| Concern | Tool |
|---|---|
| Metrics | Prometheus |
| Dashboard | Grafana |
| Logs | Loki |
| Traces | Tempo |
| Instrumentation | OpenTelemetry |

---

# Prometheus Metrics

```text
outbox_processed_total
outbox_failed_total
outbox_retry_total
outbox_queue_size
```

---

# Grafana Alert Rules

```promql
rate(outbox_failed_total[1m]) > 5
outbox_queue_size > 100
rate(outbox_processed_total[5m]) == 0
rate(outbox_retry_total[1m]) > 10
```

---

# Distributed Tracing Flow

```text
HTTP → PostgreSQL → Kafka → Consumer
```

OpenTelemetry tracing is propagated across all services.

---

# Infrastructure Stack

Docker Compose includes:

- PostgreSQL
- Kafka (KRaft mode)
- Kafka UI
- Prometheus
- Grafana
- Loki
- Tempo

---

# Docker Compose

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:15
    ports:
      - "5432:5432"

  kafka:
    image: bitnami/kafka:latest
    ports:
      - "9092:9092"

  kafka-ui:
    image: provectuslabs/kafka-ui
    ports:
      - "8080:8080"

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"

  loki:
    image: grafana/loki
    ports:
      - "3100:3100"

  tempo:
    image: grafana/tempo
    ports:
      - "3200:3200"
```

---

# CI/CD Pipeline

GitHub Actions pipeline:

- lint
- test
- build
- docker build
- docker push
- multi-service deployment pipeline

DockerHub images are built automatically for:

- order-service
- relay-worker
- payment-service
- notification-service

---

# Setup

## Start Infrastructure

```bash
docker compose up -d
```

---

## Run Order Service

```bash
cd order-service
go run cmd/api/main.go
```

---

## Run Relay Worker

```bash
cd relay-worker
go run cmd/main.go
```

---

## Run Payment Service

```bash
cd payment-service
./mvnw spring-boot:run
```

---

## Run Notification Service

```bash
cd notification-service
uvicorn app.main:app --reload
```

---

make stop-all
make run-all

# Goals

This project focuses on:

- Reliability
- Scalability
- Observability
- Distributed systems
- Event-driven architecture
- Production-grade engineering
- Polyglot microservices
- Async processing
- Fault tolerance
- Transactional consistency

---

# Production Features

Planned production-grade capabilities:

- Exponential retry backoff
- Dead Letter Queue
- Distributed tracing
- Structured logging
- Metrics aggregation
- Kafka retry topics
- DockerHub deployment
- GitHub Actions CI/CD
- Health checks
- Graceful shutdown
- Correlation IDs
- Idempotent consumers
- OpenTelemetry instrumentation

---

# Final Result

A production-style distributed backend ecosystem demonstrating:

- Clean Architecture
- Feature-based modularity
- Event-driven microservices
- Polyglot service communication
- Kafka streaming
- Observability stack
- Transactional Outbox Pattern
- CI/CD automation
- Containerized deployment
- Production engineering practices

