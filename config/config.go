package config

import (
	"fmt"
)

type Config struct {
	Debug bool `json:"debug" env:"DEBUG" default:"false"`
	//govalid:required
	//govalid:maxlength=32
	//govalid:minlength=3
	Name string    `json:"name" env:"NAME" default:"GoChat"`
	App  AppConfig `json:"app" env:", prefix=APP_"`
	DB   DbConfig  `json:"database" env:", prefix=POSTGRES_"`
}

type AppConfig struct {
	Address string `json:"addr" env:"ADDR" default:"0.0.0.0"`
	Port    uint16 `json:"port" env:"PORT" default:"8080"`
	Prefork bool   `json:"prefork" env:"PREFORK" default:"false"`
}

type DbConfig struct {
	//govalid:required
	User string `json:"user" env:"USER"`
	//govalid:required
	Password string `json:"password" env:"PASSWORD"`
	//govalid:required
	Addr string `json:"addr" env:"ADDR" default:"localhost"`
	Port uint16 `json:"port" env:"PORT" default:"5432"`
	//govalid:required
	DB string `json:"db" env:"DB" default:"chat"`
	//govalid:required
	SSLMode string `json:"ssl_mode" env:"SSLMODE" default:"disable"`
	Migrate bool   `json:"migrate" env:"MIGRATE" default:"true"`
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Address, c.App.Port)
}

func (c *DbConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Addr,
		int(c.Port),
		c.DB,
		c.SSLMode,
	)
}
