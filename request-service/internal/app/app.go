package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PanGan21/packages/auth"
	"github.com/PanGan21/packages/httpserver"
	"github.com/PanGan21/packages/logger"
	"github.com/PanGan21/packages/postgres"
	"github.com/PanGan21/request-service/config"
	requestRepository "github.com/PanGan21/request-service/internal/repository/request"
	routes "github.com/PanGan21/request-service/internal/routes"
	"github.com/PanGan21/request-service/internal/service"
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

	requestRepo := requestRepository.NewRequestRepository(*pg)
	authService := auth.NewAuthService([]byte(cfg.AuthSecret))
	requestService := service.NewRequestService(requestRepo)

	// HTTP Server
	gin.SetMode(gin.ReleaseMode)
	handler := gin.Default()

	routes.NewRouter(handler, l, authService, requestService)
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
