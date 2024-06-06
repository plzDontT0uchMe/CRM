package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env         string  `yaml:"env"`
	ApiGateway  Service `yaml:"api_gateway"`
	AuthService Service `yaml:"auth_service"`
	Secret      string  `yaml:"secret"`
	DB          DB      `yaml:"db"`
	Redis       Redis   `yaml:"redis"`
}

type Service struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Address string `yaml:"address"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"dbname"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
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
			log.Fatalf("Config file does not exist: %v", err)
		}
	}
	if err = cleanenv.ReadConfig(configFile, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}

func GetConfig() Config {
	return cfg
}
