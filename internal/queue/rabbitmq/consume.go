package rabbitmq

import (
	"context"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/observability"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func (r *RabbitMQ) Consume() (<-chan []byte, error) {
	msgChan := make(chan []byte)

	deliveries, err := r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		r.log.Error("failed to register a consumer", zap.Error(err))
		return nil, err
	}

	go func() {
		for d := range deliveries {
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), AmqpHeaderCarrier(d.Headers))
			_, span := otel.Tracer("consumer").Start(ctx, "Consume")
			if r.messageStatusCommitter != nil {
				r.messageStatusCommitter.CommitMessageStatus(r.queueName, observability.StatusMessageConsumed)
			}
			msgChan <- d.Body
			span.End()
		}
	}()

	return msgChan, nil
}
