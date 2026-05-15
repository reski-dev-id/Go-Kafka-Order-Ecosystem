package main

import (
	"log"

	"order-service/internal/bootstrap"
)

func main() {
	app, cfg, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Fatal(
		app.Start(":" + cfg.App.Port),
	)
}
