package worker

import (
	"context"
	"log"
	"time"

	"relay-worker/internal/kafka"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OutboxEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	AggregateType string
	AggregateID   uuid.UUID
	EventType     string
	Payload       datatypes.JSON
	Status        string
	RetryCount    int
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}

type RelayWorker struct {
	db       *gorm.DB
	producer *kafka.Producer
}

func NewRelayWorker(
	db *gorm.DB,
	producer *kafka.Producer,
) *RelayWorker {

	return &RelayWorker{
		db:       db,
		producer: producer,
	}
}

func (w *RelayWorker) Start(ctx context.Context) {

	log.Println("relay worker started")

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {

		case <-ctx.Done():
			log.Println("relay worker stopped")
			return

		case <-ticker.C:
			w.processEvents(ctx)
		}
	}
}

func (w *RelayWorker) processEvents(ctx context.Context) {

	var events []OutboxEvent

	err := w.db.
		Where("status = ?", "PENDING").
		Order("created_at ASC").
		Limit(10).
		Find(&events).
		Error

	if err != nil {
		log.Println(err)
		return
	}

	for _, event := range events {

		err := w.producer.Publish(
			ctx,
			event.AggregateID.String(),
			event.Payload,
		)

		if err != nil {

			log.Println("publish failed:", err)

			w.db.Model(&event).
				Update("retry_count", event.RetryCount+1)

			continue
		}

		now := time.Now()

		err = w.db.Model(&event).
			Updates(map[string]interface{}{
				"status":       "PROCESSED",
				"processed_at": now,
			}).
			Error

		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(
			"event published:",
			event.ID,
			event.EventType,
		)
	}
}
