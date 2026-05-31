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
- Relay worker publisher
- Kafka producer abstraction
- Payment service Kafka consumer
- Notification service Kafka consumer
- Event contract implementation
- Idempotent consumer handling
- Processed events tracking
- Graceful shutdown for all services

---

## Doc
📸 Screenshots:
https://github.com/reski-dev-id/Go-Kafka-Order-Ecosystem/tree/master/documentasi/screenshoot


# Architecture

```text
Client
  ↓
Order Service (Go + Echo)
  ↓
PostgreSQL
(orders + outbox_events)
  ↓
Relay Worker (Go)
  ↓
Apache Kafka
  ↓
Payment Service (Spring Boot)
  ↓
payment.completed
  ↓
Notification Service (FastAPI)
```

---

# Tech Stack

| Layer                 | Technology                   |
| --------------------- | ---------------------------- |
| Order Service         | Go 1.24                      |
| HTTP Framework        | Echo                         |
| Relay Worker          | Go 1.24                      |
| Payment Service       | Spring Boot 3                |
| Notification Service  | FastAPI                      |
| Database              | PostgreSQL 16                |
| ORM                   | GORM                         |
| Message Broker        | Apache Kafka                 |
| Kafka Client (Go)     | Sarama                       |
| Kafka Client (Java)   | Spring Kafka                 |
| Kafka Client (Python) | aiokafka                     |
| API Documentation     | Swagger / OpenAPI            |
| Monitoring            | Prometheus                   |
| Visualization         | Grafana                      |
| Kafka UI              | Provectus Kafka UI           |
| Containerization      | Docker                       |
| Orchestration         | Docker Compose               |
| Architecture          | Clean Architecture           |
| Messaging Pattern     | Event-Driven Architecture    |
| Reliability Pattern   | Transactional Outbox Pattern |
| Consumer Reliability  | Idempotent Consumer          |
| Dependency Injection  | Manual Dependency Injection  |
| CI/CD                 | GitHub Actions               |
| Artifact Registry     | DockerHub                    |

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

* PostgreSQL 16
* Apache Kafka
* Zookeeper
* Kafka UI
* Prometheus
* Grafana
* Docker Compose

Infrastructure capabilities:

* Multi-database bootstrap
* Kafka topic auto creation
* Event-driven communication
* Metrics collection
* Service monitoring
* Dashboard visualization
* Containerized local development
* Health monitoring

Infrastructure services:

| Service        | Purpose                |
| -------------- | ---------------------- |
| PostgreSQL     | Persistent storage     |
| Kafka          | Event streaming        |
| Zookeeper      | Kafka coordination     |
| Kafka UI       | Kafka topic inspection |
| Prometheus     | Metrics collection     |
| Grafana        | Metrics visualization  |
| Docker Compose | Local orchestration    |

---


# Kafka Topics

```text
order.created
payment.completed
notification.send
```

Topics are automatically created during infrastructure startup using kafka-init container.

---

# PostgreSQL Databases

The infrastructure automatically creates:

```text
order_db
payment_db
notification_db
```

---

# Graceful Shutdown

All services implement graceful shutdown handling.

Capabilities:

- Graceful HTTP server shutdown
- Kafka consumer cleanup
- Kafka producer cleanup
- Database connection cleanup
- Context cancellation
- Async task cancellation
- Consumer group leave handling
- Shutdown timeout handling

Implemented in:

- order-service
- relay-worker
- payment-service
- notification-service

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

The relay worker:

- polls pending outbox events
- publishes events to Kafka
- marks events as processed
- tracks retry count
- supports graceful shutdown
- safely closes Kafka producer

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


# Docker Compose Infrastructure

The entire ecosystem can be started locally using a single Docker Compose configuration.

Services included:

* PostgreSQL 16
* Apache Kafka
* Zookeeper
* Kafka UI
* Kafka Topic Bootstrap
* Prometheus
* Grafana


---

# Make Commands
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
# Monitoring & Observability

## Kafka UI

Available at:

```text
http://localhost:8080
```

Used for:

* Topic inspection
* Message monitoring
* Consumer group monitoring

---

## Prometheus

Available at:

```text
http://localhost:9090
```

Used for:

* Metrics collection
* Service monitoring
* PromQL queries

---

## Grafana

Available at:

```text
http://localhost:3000
```

Default credentials:

```text
username: admin
password: admin
```

Used for:

* Dashboard visualization
* Metrics monitoring
* Infrastructure observability

---


# Current Engineering Focus

This project currently emphasizes:

- Event-driven architecture
- Transactional consistency
- Reliable async messaging
- Idempotent event consumption
- Graceful shutdown handling
- Distributed backend communication
- Modular clean architecture
- Polyglot microservices
- Production-oriented infrastructure
- Fault-tolerant backend engineering

---


# Reliability Features

Implemented reliability mechanisms:

- Transactional Outbox Pattern
- Idempotent consumers
- Processed event tracking
- Graceful shutdown lifecycle
- Kafka consumer group cleanup
- Safe producer shutdown
- Async task cancellation
- Database transaction consistency
