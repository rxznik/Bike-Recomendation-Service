package mockproducer

import (
	"context"
	"time"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/queue"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type MockProducer struct {
	producer queue.Producer
	log      *zap.Logger
}

func New(logger *zap.Logger, producer queue.Producer) *MockProducer {
	return &MockProducer{
		producer: producer,
		log:      logger,
	}
}

func (m *MockProducer) Run() error {

	for {
		for _, msg := range mockMessages {

			msgJSON, err := json.Marshal(msg)
			if err != nil {
				m.log.Error("failed to marshal message", zap.Error(err))
				return err
			}

			if err := m.producer.Produce(context.Background(), msgJSON); err != nil {
				m.log.Error("failed to produce", zap.Error(err))
				return err
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (m *MockProducer) Stop() error {
	return m.producer.Close()
}
