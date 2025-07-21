package adapters

import (
	"context"
	"encoding/json"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.uber.org/zap"
)

type WorkStatusCommitter interface {
	CommitWorkStatus(status string)
}

type Producer interface {
	Produce(ctx context.Context, message []byte) error
	Close() error
}

type Adapter struct {
	producer Producer
	log      *zap.Logger
}

func New(logger *zap.Logger, producer Producer) *Adapter {
	return &Adapter{
		producer: producer,
		log:      logger,
	}
}

func (a *Adapter) Produce(ctx context.Context, message models.RecommendationsMessage) error {

	messageJSON, err := json.Marshal(message)
	if err != nil {
		a.log.Error("failed to marshal message", zap.Error(err))
		return err
	}
	a.log.Info("message transmitted to producer", zap.String("message", string(messageJSON)))
	return a.producer.Produce(ctx, messageJSON)
}

func (a *Adapter) Close() error {
	return a.producer.Close()
}
