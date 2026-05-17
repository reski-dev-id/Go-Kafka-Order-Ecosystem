package consumer

import (
	"context"
	"encoding/json"
	"log"

	"order-service/features/order/dto"
	"order-service/features/order/repository"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type PaymentCompletedConsumer struct {
	orderRepo repository.OrderRepository
}

func NewPaymentCompletedConsumer(
	orderRepo repository.OrderRepository,
) *PaymentCompletedConsumer {
	return &PaymentCompletedConsumer{
		orderRepo: orderRepo,
	}
}

func (c *PaymentCompletedConsumer) Setup(
	sarama.ConsumerGroupSession,
) error {
	return nil
}

func (c *PaymentCompletedConsumer) Cleanup(
	sarama.ConsumerGroupSession,
) error {
	return nil
}

func (c *PaymentCompletedConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for message := range claim.Messages() {

		var event dto.PaymentCompletedEvent

		err := json.Unmarshal(
			message.Value,
			&event,
		)

		if err != nil {
			log.Println("failed unmarshal:", err)
			continue
		}

		orderID, err := uuid.Parse(event.OrderID)
		if err != nil {
			log.Println("failed parse uuid:", err)
			continue
		}

		err = c.orderRepo.UpdateStatusByID(
			context.Background(),
			orderID,
			event.Status,
		)

		if err != nil {
			log.Println("failed update order:", err)
			continue
		}

		log.Println(
			"order updated to PAID:",
			event.OrderID,
		)

		session.MarkMessage(message, "")
	}

	return nil
}
