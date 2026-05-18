package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "order-service/docs"

	"order-service/features/order/consumer"
	"order-service/features/order/repository"
	"order-service/internal/bootstrap"

	"github.com/IBM/sarama"
)

// @title Order Service API
// @version 1.0
// @description Order Service API Documentation
// @host localhost:8080
// @BasePath /
func main() {

	app, db, cfg, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(db)

	kafkaConfig := sarama.NewConfig()

	kafkaConfig.Version = sarama.V3_6_0_0

	kafkaConfig.Consumer.Group.Rebalance.Strategy =
		sarama.BalanceStrategyRange

	kafkaConfig.Consumer.Offsets.Initial =
		sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{"localhost:9094"},
		"order-group",
		kafkaConfig,
	)

	if err != nil {
		log.Fatal(err)
	}

	paymentCompletedConsumer :=
		consumer.NewPaymentCompletedConsumer(
			orderRepo,
		)

	go func() {

		for {

			err := consumerGroup.Consume(
				context.Background(),
				[]string{"payment.completed"},
				paymentCompletedConsumer,
			)

			if err != nil {
				log.Println(
					"consumer error:",
					err,
				)
			}
		}
	}()

	go func() {

		err := app.Start(
			":" + cfg.App.Port,
		)

		if err != nil &&
			err != http.ErrServerClosed {

			app.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	err = consumerGroup.Close()
	if err != nil {
		log.Println(
			"failed close kafka consumer:",
			err,
		)
	}

	sqlDB, err := db.DB()
	if err == nil {

		err = sqlDB.Close()

		if err != nil {
			log.Println(
				"failed close database:",
				err,
			)
		}
	}

	err = app.Shutdown(ctx)
	if err != nil {
		log.Fatal(
			"server forced shutdown:",
			err,
		)
	}

	log.Println("server exited properly")
}
