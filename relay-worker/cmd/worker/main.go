package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"relay-worker/config"
	"relay-worker/internal/database"
	"relay-worker/internal/kafka"
	"relay-worker/internal/worker"
)

func main() {

	cfg := config.LoadConfig()

	db := database.NewPostgres(cfg)

	producer := kafka.NewProducer(
		cfg.KafkaBroker,
		cfg.KafkaTopicOrderCreated,
	)

	relayWorker := worker.NewRelayWorker(
		db,
		producer,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go relayWorker.Start(ctx)

	stopChan := make(chan os.Signal, 1)

	signal.Notify(
		stopChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stopChan

	cancel()
}
