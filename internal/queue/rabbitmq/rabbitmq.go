package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type AmqpHeaderCarrier amqp.Table

func (a AmqpHeaderCarrier) Set(key, value string) {
	a[key] = value
}

func (a AmqpHeaderCarrier) Get(key string) string {
	return a[key].(string)
}

func (a AmqpHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	return keys
}

type MessageStatusCommitter interface {
	CommitMessageStatus(queue, status string)
}

type RabbitMQ struct {
	queueName              string
	conn                   *amqp.Connection
	channel                *amqp.Channel
	log                    *zap.Logger
	messageStatusCommitter MessageStatusCommitter
}

type QueueOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func NewConnection(logger *zap.Logger, rabbitMQDSN string) *RabbitMQ {
	conn, err := amqp.Dial(rabbitMQDSN)
	if err != nil {
		logger.Fatal("failed to connect to RabbitMQ", zap.Error(err))
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("failed to open a channel", zap.Error(err))
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		log:     logger,
	}
}

func (r *RabbitMQ) SetupMessageStatusCommitter(messageStatusComitter MessageStatusCommitter) {
	r.messageStatusCommitter = messageStatusComitter
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

func (r *RabbitMQ) DeclareQueue(opts QueueOptions) error {
	_, err := r.channel.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
	r.queueName = opts.Name
	return err
}
