package config

import (
	"log"
	"os"
	"time"

	envloader "github.com/devprod-tech/webike_recomendations-Vitalya/internal/utils/env-loader"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Env               string        `yaml:"env" env-default:"dev"`
	Consumer          Queue         `yaml:"consumer" env-required:"true"`
	Producer          Queue         `yaml:"producer" env-required:"true"`
	BrokerType        string        `yaml:"broker_type" env-required:"true" env-default:"rabbitmq"`
	HttpClientTimeout time.Duration `yaml:"http_client_timeout" env-required:"true" env-default:"5s"`
	GoogleAPIKey      string        `env:"GOOGLE_API_KEY" env-required:"true"`
	SentryDSN         string        `env:"SENTRY_DSN" env-required:"true"`
	BrokerDSN         string
}

type Queue struct {
	QueueName string `yaml:"queue_name" env-required:"true"`
}

func MustLoad() *Config {

	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		zap.L().Fatal("config file not found", zap.Error(err))
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	cfg.loadBrokerDSN()

	return &cfg
}

func (c *Config) loadBrokerDSN() {
	switch c.BrokerType {
	case "rabbitmq":
		c.BrokerDSN = envloader.LoadRabbitMQDSN()
	}
}

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/dev/config.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %v", configPath, err)
	}
	return configPath
}
