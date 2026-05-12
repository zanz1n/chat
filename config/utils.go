package config

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/creasty/defaults"
	"github.com/sethvargo/go-envconfig"
)

var (
	config        = flag.String("config", "/etc/chat/config.yml", "Config path")
	configDebug   = flag.Bool("debug", false, "Enables or disables debug logging")
	configMigrate = flag.Bool("migrate", false, "Executes migrations before running the server")
	configPort    = flag.Int("port", -1, "The port in which the app will run")
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if config == nil || *config == "" {
		config = new(os.Getenv("CONFIG_FILE"))
	}
}

var instance = sync.OnceValue(func() *Config {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("Failed to init config: %s\n", err)
	}
	return cfg
})

func Get() *Config {
	return instance()
}

func setFromEnv(c *Config) {
	if configDebug != nil && *configDebug {
		c.Debug = *configDebug
	}
	if configMigrate != nil && *configMigrate {
		c.DB.Migrate = true
	}
	if configPort != nil && *configPort != -1 {
		c.App.Address = "0.0.0.0"
		c.App.Port = uint16(*configPort)
	}
}

func initConfig() (*Config, error) {
	var cfg Config

	var (
		err   error
		write func(any) error
	)
	if *config == "env" {
		err = envconfig.Process(context.Background(), &cfg)
		write = func(any) error { return nil }
	} else {
		buf, err := os.ReadFile(*config)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("open config file: %w", err)
			}
		}

		if buf != nil {
			err = json.Unmarshal(buf, &cfg)
		}
		write = func(in any) error {
			out, err := json.MarshalIndent(in, "", "  ")
			if err != nil {
				return err
			}
			out = append(out, '\n')
			return os.WriteFile(*config, out, 0666)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("parse configuration: %w", err)
	}

	defaults.MustSet(&cfg)
	original := cfg
	setFromEnv(&cfg)

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	if err = cfg.DB.Validate(); err != nil {
		return nil, fmt.Errorf("invalid db config: %w", err)
	}

	if err = write(original); err != nil {
		slog.Warn("Config: Failed to write config", "error", err)
	}

	return &cfg, nil
}
