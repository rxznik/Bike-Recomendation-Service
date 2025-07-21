package envloader

import (
	"fmt"
	"os"
)

func LoadRabbitMQDSN() string {
	if rabbitMQDSN := os.Getenv("RABBITMQ_DSN"); rabbitMQDSN != "" {
		return rabbitMQDSN
	}

	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)
}
