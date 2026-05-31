package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"relay-worker/config"
	"relay-worker/internal/database"
	"relay-worker/internal/kafka"
	"relay-worker/internal/metrics"
	"relay-worker/internal/worker"
)

func main() {

	metrics.Init()

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

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		relayWorker.Start(ctx)
	}()

	stopChan := make(chan os.Signal, 1)

	signal.Notify(
		stopChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stopChan

	log.Println("shutdown signal received")

	cancel()

	log.Println("waiting relay worker to stop")

	wg.Wait()

	log.Println("closing kafka producer")

	err := producer.Close()

	if err != nil {
		log.Printf(
			"failed close producer: %v",
			err,
		)
	}

	log.Println("relay worker shutdown complete")
}
