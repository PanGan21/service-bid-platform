package main

import (
	"fmt"
	"log"

	"github.com/PanGan21/request-service/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	fmt.Println(cfg.App.Name)

	// Run
}
