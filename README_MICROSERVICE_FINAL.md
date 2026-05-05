# Go Kafka Order Ecosystem

Sistem **Order Processing berbasis Event-Driven Architecture** menggunakan **Transaction Outbox Pattern** untuk menjamin konsistensi antara database dan Kafka.

Built with **Golang**, **PostgreSQL**, **Apache Kafka**, serta dilengkapi **Prometheus, Grafana, dan OpenTelemetry** untuk observability.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22+ |
| HTTP Framework | Gin |
| Database | PostgreSQL + sqlx |
| Message Broker | Apache Kafka |
| Metrics | Prometheus |
| Dashboard | Grafana |
| Tracing | OpenTelemetry + Jaeger |
| DI Container | Wire (optional) |
| Containerization | Docker + Docker Compose |

---

## Architecture Overview

```
Client
  ↓
Order Service (HTTP)
  ↓
PostgreSQL (orders + outbox)
  ↓
Relay Worker (Outbox Processor)
  ↓
Kafka
  ↓ ↓
Payment Service      Notification Service
```

---

## Clean Architecture Layer (Per Service)

```
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

---

## Struktur Project (Multi-Service)

```
go-kafka-order-ecosystem/
│
├── order-service/
├── relay-worker/
├── payment-service/
├── notification-service/
│
├── infra/
│   ├── docker-compose.yml
│   ├── prometheus.yml
│   └── grafana-dashboard.json
│
└── README.md
```

---

## Database Schema

```sql
CREATE TABLE orders (
  id UUID PRIMARY KEY,
  item_name VARCHAR(255),
  amount DECIMAL,
  status VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE outbox (
  id BIGSERIAL PRIMARY KEY,
  aggregate_id UUID,
  topic VARCHAR(100),
  payload JSONB,
  processed BOOLEAN DEFAULT FALSE,
  retry_count INT DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Kafka Topics

- order-events
- payment-events
- notification-events
- order-events-dlq

---

## Event Flow

1. Order Service → insert order + outbox
2. Relay Worker → publish Kafka
3. Payment Service → process payment
4. Notification Service → send notification

---

## Prometheus Metrics

- outbox_processed_total
- outbox_failed_total
- outbox_retry_total
- outbox_queue_size

---

## Grafana Alert Rules

```
rate(outbox_failed_total[1m]) > 5
outbox_queue_size > 100
rate(outbox_processed_total[5m]) == 0
rate(outbox_retry_total[1m]) > 10
```

---

## OpenTelemetry Tracing

Flow:
HTTP → DB → Kafka → Consumer

---

## Docker Compose (Infra)

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:15
    ports:
      - "5432:5432"

  kafka:
    image: bitnami/kafka
    ports:
      - "9092:9092"

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"

  jaeger:
    image: jaegertracing/all-in-one
    ports:
      - "16686:16686"
```

---

## Setup

```bash
docker-compose up -d

cd order-service && go run cmd/api/main.go
cd relay-worker && go run cmd/main.go
cd payment-service && go run cmd/main.go
cd notification-service && go run cmd/main.go
```

---

## Result

- Reliable
- Scalable
- Observable
- Production-ready
