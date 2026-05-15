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

	"order-service/internal/bootstrap"
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

	go func() {
		err := app.Start(":" + cfg.App.Port)
		if err != nil && err != http.ErrServerClosed {
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

	sqlDB, err := db.DB()
	if err == nil {
		err = sqlDB.Close()
		if err != nil {
			log.Println("failed close database:", err)
		}
	}

	err = app.Shutdown(ctx)
	if err != nil {
		log.Fatal("server forced shutdown:", err)
	}

	log.Println("server exited properly")
}
