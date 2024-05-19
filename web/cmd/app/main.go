package main

import (
	"log"

	"github.com/abielalejandro/web/config"
	"github.com/abielalejandro/web/internals/app"
)

func main() {
	log.Println("Starting app")

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app := app.NewApp(cfg)
	app.Run()
}
