package config

type Config struct {
	Port string `yaml:"port" env:"PORT" env-default:"3000"`
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
}

var cfg Config

func GetConfig() Config {
	return cfg
}

func SetConfig(config *Config) {
	cfg = *config
}
