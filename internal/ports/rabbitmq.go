package ports

import (
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/queue/rabbitmq"
	envloader "github.com/devprod-tech/webike_recomendations-Vitalya/internal/utils/env-loader"
	"go.uber.org/zap"
)

type RabbitMQPortOptions struct {
	RabbitMQDSN            string
	QueueOptions           *rabbitmq.QueueOptions
	MessageStatusCommitter rabbitmq.MessageStatusCommitter
}

func NewRabbitMQ(logger *zap.Logger, opts *RabbitMQPortOptions) *Port {
	if opts == nil {
		opts = &RabbitMQPortOptions{
			RabbitMQDSN:  envloader.LoadRabbitMQDSN(),
			QueueOptions: &rabbitmq.QueueOptions{Durable: true},
		}
	}
	consumer := rabbitmq.NewConnection(logger, opts.RabbitMQDSN)

	consumer.SetupMessageStatusCommitter(opts.MessageStatusCommitter)

	if err := consumer.DeclareQueue(*opts.QueueOptions); err != nil {
		logger.Fatal("failed to declare a queue", zap.Error(err))
	}
	logger.Info("consumer created", zap.String("queue", opts.QueueOptions.Name))

	return &Port{
		outputCh: make(chan models.AnalyticsMessage),
		consumer: consumer,
		log:      logger,
	}
}
