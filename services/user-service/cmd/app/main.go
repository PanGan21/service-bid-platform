package main

import (
	"log"

	"github.com/PanGan21/user-service/config"
	"github.com/PanGan21/user-service/internal/app"
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
