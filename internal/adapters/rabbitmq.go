package adapters

import (
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/queue/rabbitmq"
	envloader "github.com/devprod-tech/webike_recomendations-Vitalya/internal/utils/env-loader"
	"go.uber.org/zap"
)

type RabbitMQAdapterOptions struct {
	RabbitMQDSN            string
	QueueOptions           *rabbitmq.QueueOptions
	MessageStatusCommitter rabbitmq.MessageStatusCommitter
}

func NewRabbitMQ(logger *zap.Logger, opts *RabbitMQAdapterOptions) *Adapter {
	if opts == nil {
		opts = &RabbitMQAdapterOptions{
			RabbitMQDSN:  envloader.LoadRabbitMQDSN(),
			QueueOptions: &rabbitmq.QueueOptions{Durable: true},
		}
	}
	producer := rabbitmq.NewConnection(logger, opts.RabbitMQDSN)

	producer.SetupMessageStatusCommitter(opts.MessageStatusCommitter)

	if err := producer.DeclareQueue(*opts.QueueOptions); err != nil {
		logger.Fatal("failed to declare a queue", zap.Error(err))
	}
	logger.Info("producer created", zap.String("queue", opts.QueueOptions.Name))

	return &Adapter{
		producer: producer,
		log:      logger,
	}
}
