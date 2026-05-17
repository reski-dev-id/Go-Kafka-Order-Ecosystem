package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(
	broker string,
	topic string,
) *Producer {

	writer := &kafkago.Writer{
		Addr:     kafkago.TCP(broker),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	key string,
	value []byte,
) error {

	return p.writer.WriteMessages(
		ctx,
		kafkago.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
}
