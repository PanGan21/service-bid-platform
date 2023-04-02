package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Postgres `yaml:"postgres"`
		Kafka    `yaml:"kafka"`
	}

	App struct {
		Name string `env-required:"true" yaml:"name"    env:"APP_NAME"`
	}

	// HTTP -.
	HTTP struct {
		Port          string   `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		SessionSecret string   `env-required:"true" yaml:"session_secret" env:"SESSION_SECRET"`
		AuthSecret    string   `env-required:"true" yaml:"auth_secret" env:"AUTH_SECRET"`
		CorsOrigins   []string `env-required:"true" yaml:"cors_origins" env:"CORS_ORIGINS"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// Postgres -.
	Postgres struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"url" env:"PG_URL"`
	}

	// Kafka -.
	Kafka struct {
		Retries int    `env-required:"true" yaml:"retries" env:"KAFKA_RETRIES"`
		URL     string `env-required:"true" yaml:"url" env:"KAFKA_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
