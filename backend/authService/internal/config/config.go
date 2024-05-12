package config

type Config struct {
	Port       string `yaml:"port" env:"PORT" env-default:"3001"`
	Host       string `yaml:"host" env:"HOST" env-default:"localhost"`
	DBHost     string `yaml:"db_host" env:"DB_HOST" env-default:"localhost"`
	DBPort     string `yaml:"db_port" env:"DB_PORT" env-default:"5432"`
	DBUser     string `yaml:"db_user" env:"DB_USER" env-default:"postgres"`
	DBPassword string `yaml:"db_password" env:"DB_PASSWORD" env-default:"0000"`
	DBName     string `yaml:"db_name" env:"DB_NAME" env-default:"crm"`
	Secret     string `yaml:"secret" env:"SECRET" env-default:"$ecr3t"`
}

var cfg Config

func GetConfig() Config {
	return cfg
}

func SetConfig(config *Config) {
	cfg = *config
}
