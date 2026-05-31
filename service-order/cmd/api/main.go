package main

import (
	"context"
	"errors"
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
	"order-service/internal/metrics"

	"github.com/IBM/sarama"
)

// @title Order Service API
// @version 1.0
// @description Order Service API Documentation
// @host localhost:8081
// @BasePath /
func main() {

	metrics.Init()

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

	consumerCtx, consumerCancel :=
		context.WithCancel(context.Background())

	go func() {

		for {

			select {

			case <-consumerCtx.Done():

				log.Println(
					"payment consumer stopped",
				)

				return

			default:

				err := consumerGroup.Consume(
					consumerCtx,
					[]string{"payment.completed"},
					paymentCompletedConsumer,
				)

				if err != nil {

					if errors.Is(
						err,
						context.Canceled,
					) {
						return
					}

					log.Println(
						"consumer error:",
						err,
					)

					time.Sleep(
						2 * time.Second,
					)
				}
			}
		}
	}()

	go func() {

		log.Println(
			"order service started on port:",
			cfg.App.Port,
		)

		err := app.Start(
			":" + cfg.App.Port,
		)

		if err != nil &&
			!errors.Is(
				err,
				http.ErrServerClosed,
			) {

			app.Logger.Fatal(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)

	signal.Notify(
		stopChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stopChan

	log.Println(
		"shutdown signal received",
	)

	consumerCancel()

	shutdownCtx, cancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

	defer cancel()

	log.Println(
		"shutting down http server",
	)

	err = app.Shutdown(shutdownCtx)

	if err != nil {

		log.Println(
			"failed shutdown http server:",
			err,
		)
	}

	log.Println(
		"closing kafka consumer group",
	)

	err = consumerGroup.Close()

	if err != nil {

		log.Println(
			"failed close kafka consumer:",
			err,
		)
	}

	log.Println(
		"closing database connection",
	)

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

	log.Println(
		"order service shutdown complete",
	)
}
