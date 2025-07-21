package rabbitmq

import (
	"context"
	"time"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/observability"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (r *RabbitMQ) Produce(ctx context.Context, message []byte) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("producer").Start(ctx, "Produce")
	defer span.End()

	headers := make(AmqpHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, headers)

	if r.messageStatusCommitter != nil {
		defer func() {
			if err != nil {
				r.messageStatusCommitter.CommitMessageStatus(r.queueName, observability.StatusMessageProduceFailed)
				r.log.Error("failed to publish a message", zap.Error(err))
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			} else {
				r.messageStatusCommitter.CommitMessageStatus(r.queueName, observability.StatusMessageProduced)
			}
		}()
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			Headers:     amqp.Table(headers),
			ContentType: "application/json",
			Body:        message,
		},
	)

	return err
}
