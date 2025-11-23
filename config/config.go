package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"url"      env:"PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("⚠️  .env file not found, trying .env.example...")
		err = godotenv.Load(".env.example")
		if err != nil {
			fmt.Printf("⚠️  .env.example file not found: %v\n", err)
			fmt.Println("ℹ️  Continuing with environment variables and config.yml...")
		} else {
			fmt.Println("✅ Loaded configuration from .env.example")
		}
	} else {
		fmt.Println("✅ Loaded configuration from .env")
	}

	cfg := &Config{}

	err = cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
