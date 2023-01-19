package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PanGan21/pkg/auth"
	"github.com/PanGan21/pkg/httpserver"
	"github.com/PanGan21/pkg/logger"
	"github.com/PanGan21/pkg/postgres"
	"github.com/PanGan21/user-service/config"
	userRepository "github.com/PanGan21/user-service/internal/repository/user"
	routes "github.com/PanGan21/user-service/internal/routes/http"
	"github.com/PanGan21/user-service/internal/service"
	"github.com/gin-gonic/contrib/sessions"
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

	// Session store
	store, err := sessions.NewRedisStore(10, "tcp", cfg.Redis.Url, "", []byte(cfg.SessionSecret))
	if err != nil {
		l.Fatal(err)
	}

	userRepo := userRepository.NewUserRepository(*pg)
	userService := service.NewUserService(userRepo, cfg.User.PasswordSalt)
	authService := auth.NewAuthService([]byte(cfg.AuthSecret))

	// HTTP Server
	gin.SetMode(gin.ReleaseMode)
	handler := gin.Default()

	routes.NewRouter(handler, l, store, userService, authService)
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
