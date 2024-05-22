package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	AuthService HTTPServer `yaml:"auth_service"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:3001"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
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
