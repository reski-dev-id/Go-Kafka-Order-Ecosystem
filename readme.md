# Go Kafka Order Ecosystem

Production-style distributed backend ecosystem implementing Event-Driven Architecture using the Transactional Outbox Pattern for reliable asynchronous communication between services.

This project focuses on:

- Microservice architecture
- Event-driven systems
- Polyglot services
- Clean Architecture
- Distributed messaging with Kafka
- Transactional consistency
- Async processing
- Dockerized infrastructure
- Production-oriented backend engineering

---

# Current Project Status

## Implemented

- Dockerized infrastructure
- PostgreSQL 16
- Apache Kafka
- Zookeeper
- Kafka UI
- Kafka topic auto bootstrap
- Multi-database bootstrap
- Order CRUD API
- Clean Architecture
- Feature-based structure
- DTO validation
- Centralized validation middleware
- Standardized API response
- Transactional Outbox Pattern
- PostgreSQL outbox persistence

---

## In Progress

- Relay worker
- Kafka publisher abstraction
- Payment service Kafka consumer
- Notification service Kafka consumer
- Retry mechanism
- Idempotent consumer handling

---

## Planned

- OpenTelemetry tracing
- Prometheus metrics
- Grafana dashboards
- Loki logging
- Tempo tracing
- Dead Letter Queue (DLQ)
- Kafka retry topics
- GitHub Actions CI/CD
- DockerHub deployment
- Flyway migration management
- Google Wire dependency injection
- Distributed tracing propagation
- Structured logging

---

# Architecture

```text
Client
  ↓
Order Service (Go)
  ↓
PostgreSQL
(orders + outbox)
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
| Order Service | Go 1.22 |
| HTTP Framework | Echo |
| Relay Worker | Go |
| Payment Service | Spring Boot |
| Notification Service | FastAPI |
| Database | PostgreSQL 16 |
| SQL Library | sqlx |
| Message Broker | Apache Kafka |
| Kafka UI | Provectus Kafka UI |
| Containerization | Docker + Docker Compose |
| Architecture | Clean Architecture |
| Dependency Injection | Manual DI (planned Wire) |

---

# Microservices

| Service | Responsibility |
|---|---|
| order-service | Create orders + store outbox events |
| relay-worker | Publish outbox events to Kafka |
| payment-service | Consume order events + process payment |
| notification-service | Consume payment events + send notifications |

---

# Clean Architecture

Each service follows Clean Architecture principles.

```text
┌────────────────────────────┐
│       Delivery Layer       │
├────────────────────────────┤
│       Usecase Layer        │
├────────────────────────────┤
│      Repository Layer      │
├────────────────────────────┤
│        Domain Layer        │
└────────────────────────────┘
```

Rules:

- Domain layer remains framework-independent
- Usecase layer contains business logic
- Infrastructure implements interfaces
- Dependencies flow inward only

---

# Feature-Based Structure

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

- Better modularity
- Easier scalability
- Reduced coupling
- Clear ownership boundaries

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
│   ├── postgres/
│   ├── kafka/
│   ├── middleware/
│   ├── response/
│   └── errors/
│
├── order-service/
│   ├── cmd/
│   ├── internal/
│   └── migrations/
│
├── relay-worker/
│
├── payment-service/
│
├── notification-service/
│
├── postgres-init/
│
├── infra/
│
├── .github/
│
└── README.md
```

---

# Infrastructure Stack

Current infrastructure:

- PostgreSQL 16
- Apache Kafka
- Zookeeper
- Kafka UI
- Docker Compose

---

# Kafka Topics

```text
order.created
payment.completed
notification.send
```

Topics are automatically bootstrapped during container startup.

---

# PostgreSQL Databases

The infrastructure automatically creates:

```text
order_db
payment_db
notification_db
```

---

# Database Schema

## Orders Table

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_name VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Outbox Events Table

```sql
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);
```

---

# Transactional Outbox Pattern

The system uses the Transactional Outbox Pattern to guarantee consistency between PostgreSQL transactions and Kafka event publishing.

Flow:

```text
HTTP Request
  ↓
Create Order
  ↓
Insert Outbox Event
  ↓
Commit Transaction
  ↓
Relay Worker Polling
  ↓
Kafka Publish
```

Benefits:

- Prevents dual-write problems
- Guarantees event persistence
- Reliable asynchronous messaging
- Better fault tolerance

---

# Relay Worker Flow

The relay worker will:

- poll unprocessed outbox events
- publish events to Kafka
- mark events as processed
- retry failed events

Polling query:

```sql
SELECT *
FROM outbox_events
WHERE status = 'PENDING'
ORDER BY created_at
LIMIT 10
FOR UPDATE SKIP LOCKED;
```

---

# Event Contract

Example event payload:

```json
{
  "event_id": "uuid",
  "event_type": "order.created",
  "aggregate_id": "uuid",
  "occurred_at": "timestamp",
  "payload": {
    "order_id": "uuid",
    "customer_name": "reski",
    "amount": 2500
  }
}
```

---

# Docker Compose Infrastructure

Services included:

- PostgreSQL
- Kafka
- Zookeeper
- Kafka UI
- Kafka topic bootstrap

---

# Make Commands

## Start Infrastructure

```bash
make infra-up
```

---

## Stop Infrastructure

```bash
make infra-down
```

---

## Reset Infrastructure

```bash
make infra-reset
```

---

## Run Order Service

```bash
make run-order
```

---

## Run Relay Worker

```bash
make run-relay
```

---

## Run Payment Service

```bash
make run-payment
```

---

## Run Notification Service

```bash
make run-notification
```

---

## Run All Services

```bash
make run-all
```

---

## Stop All Services

```bash
make stop-all
```

---

# Manual Commands

## Run Infrastructure

```bash
docker compose up -d
```

---

## Stop Infrastructure

```bash
docker compose down
```

---

## Reset Infrastructure

```bash
docker compose down -v
```

---

# Kafka UI

Available at:

```text
http://localhost:8080
```

---

# Current Engineering Focus

This project currently emphasizes:

- Reliable event publishing
- Async processing
- Transaction consistency
- Modular architecture
- Distributed backend communication
- Production-style infrastructure
- Polyglot microservices

---

# Planned Production Features

Future production-grade capabilities:

- Dead Letter Queue
- Retry topics
- Exponential retry backoff
- Distributed tracing
- Metrics aggregation
- Structured logging
- Correlation IDs
- OpenTelemetry instrumentation
- Health checks
- Graceful shutdown
- CI/CD automation
- DockerHub deployment
- Observability stack

---

# Final Goal

A production-style distributed backend ecosystem demonstrating:

- Event-driven microservices
- Transactional Outbox Pattern
- Clean Architecture
- Kafka streaming
- Polyglot services
- Distributed systems engineering
- Containerized infrastructure
- Production-oriented backend practices