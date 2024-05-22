package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	App        *App        `env:",prefix=APP_"`
	Monitoring *Monitoring `env:",prefix=MONITORING_"`
	Http       *Http       `env:",prefix=HTTP_"`
}

type App struct {
	Port string `env:"PORT"`
}

type Monitoring struct {
	Port        string        `env:"PORT"`
	ReadTimeout time.Duration `env:"READ_TIMEOUT"`
}

type Http struct {
	ReadTimeout time.Duration `env:"READ_TIMEOUT"`
}

func New(ctx context.Context) (*Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("process: %w", err)
	}
	return &cfg, nil
}
