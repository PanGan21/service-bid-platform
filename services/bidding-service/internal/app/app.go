package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PanGan21/bidding-service/config"
	requestRepository "github.com/PanGan21/bidding-service/internal/repository/request"
	"github.com/PanGan21/bidding-service/internal/routes/events"
	routes "github.com/PanGan21/bidding-service/internal/routes/http"
	"github.com/PanGan21/bidding-service/internal/service"
	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/httpserver"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/messaging"
	"github.com/PanGan21/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	var err error

	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Events
	sub := messaging.NewSubscriber(cfg.Kafka.URL, cfg.App.Name)

	requestRepo := requestRepository.NewRequestRepository(*pg)
	authService := auth.NewAuthService([]byte(cfg.AuthSecret))
	requestService := service.NewRequestService(requestRepo)

	events.NewEventsClient(sub, l, requestService)

	// HTTP Server
	gin.SetMode(gin.ReleaseMode)
	handler := gin.Default()

	// sub.Subscribe("request-created", handlePayload)
	routes.NewRouter(handler, l, authService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
