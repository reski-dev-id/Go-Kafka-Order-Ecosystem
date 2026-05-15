package main

import (
	"log"

	_ "order-service/docs"

	"order-service/internal/bootstrap"
)

// @title Order Service API
// @version 1.0
// @description Order Service API Documentation
// @host localhost:8080
// @BasePath /
func main() {
	app, cfg, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Fatal(
		app.Start(":" + cfg.App.Port),
	)
}
