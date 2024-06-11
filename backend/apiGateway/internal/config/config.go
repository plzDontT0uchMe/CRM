package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env             string  `yaml:"env"`
	ApiGateway      Service `yaml:"api_gateway"`
	AuthService     Service `yaml:"auth_service"`
	UsersService    Service `yaml:"users_service"`
	StorageService  Service `yaml:"storage_service"`
	SubsService     Service `yaml:"subs_service"`
	TrainingService Service `yaml:"training_service"`
}

type Service struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Address string `yaml:"address"`
}

var cfg Config

func init() {
	configFile, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory for config: %v", err)
	}
	configFile += "\\..\\..\\internal\\config\\config.yaml"
	if _, err = os.Stat(configFile); err != nil {
		configFile, _ = os.Getwd()
		configFile += "/config.yaml"
		if _, err = os.Stat(configFile); err != nil {
			log.Fatalf("Config file.jpg does not exist: %v", err)
		}
	}
	if err = cleanenv.ReadConfig(configFile, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}

func GetConfig() Config {
	return cfg
}
