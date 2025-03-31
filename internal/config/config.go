package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		Env string `yaml:"env" env-default:"local"`
		DB
		HTTPServer `yaml:"http_server"`
	}
	HTTPServer struct {
		Address     string        `yaml:"address" env-default:"localhost:8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
		User string `yaml:"user" env-required:"true"`
		Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
	}
	DB struct {
		DataSourceName string `env:"DB_URL" env-required:"true"`
	}
)

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load env variables: %s", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("MustLoad: unable to load cfg: %s", err)
	}
	return &cfg
}
