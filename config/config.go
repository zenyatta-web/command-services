package config

import (
	"fmt"

	"go-micro.dev/v4/config"            //permite manejar la configuracion.
	"go-micro.dev/v4/config/source/env" //permite cargar la configuracion desde variables de entorno.
)

type Config struct {
	Port    int
	Tracing TraicingConfig //rastreo
}

type TraicingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

var cfg *Config = &Config{
	Port: 50052,
}

func Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Tracing() TraicingConfig {
	return cfg.Tracing
}

func Load() error {
	configor, err := config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		return fmt.Errorf("configor.New: %w", err)
	}
	if err := configor.Load(); err != nil {
		return fmt.Errorf("configor.Load: %w", err)
	}
	if err := configor.Scan(cfg); err != nil {
		return fmt.Errorf("configor.Scan: %w", err)
	}
	return nil
}
