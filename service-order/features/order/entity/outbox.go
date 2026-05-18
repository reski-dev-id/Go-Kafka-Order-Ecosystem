package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
