package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Secret     string     `yaml:"secret" env-default:"$ecr3t"`
	DB         DB         `yaml:"db"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:3001"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
}

type DB struct {
	DBHost     string `yaml:"host" env-default:"localhost"`
	DBPort     int    `yaml:"port" env-default:"5432"`
	DBUser     string `yaml:"user" env-default:"postgres"`
	DBPassword string `yaml:"password" env-default:"0000"`
	DBName     string `yaml:"dbname" env-default:"crm-authService"`
}

var cfg Config

func init() {
	configFile, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory for config: %v", err)
	}
	configFile += "\\..\\..\\internal\\config\\config.yaml"
	if _, err = os.Stat(configFile); err != nil {
		log.Fatalf("Config file does not exist: %v", err)
	}
	if err = cleanenv.ReadConfig(configFile, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}

func GetConfig() Config {
	return cfg
}
