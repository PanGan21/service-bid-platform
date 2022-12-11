package main

import (
	"log"

	"github.com/PanGan21/bidding-service/config"
	"github.com/PanGan21/bidding-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
