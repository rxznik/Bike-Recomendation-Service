package ports

import (
	"encoding/json"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.uber.org/zap"
)

type Consumer interface {
	Consume() (<-chan []byte, error)
	Close() error
}

type Port struct {
	outputCh chan models.AnalyticsMessage
	consumer Consumer
	log      *zap.Logger
}

func New(logger *zap.Logger, consumer Consumer) *Port {
	return &Port{
		outputCh: make(chan models.AnalyticsMessage),
		consumer: consumer,
		log:      logger,
	}
}

func (p *Port) Open() (<-chan models.AnalyticsMessage, error) {

	msgChan, err := p.consumer.Consume()
	if err != nil {
		p.log.Error("failed to consume", zap.Error(err))
		return nil, err
	}

	go work(p.log, msgChan, p.outputCh)

	return p.outputCh, nil
}

func work(
	logger *zap.Logger,
	input <-chan []byte,
	output chan<- models.AnalyticsMessage,
) {
	for msg := range input {
		logger.Info("message received from consumer", zap.String("message", string(msg)))

		var analyticsMsg models.AnalyticsMessage

		if err := json.Unmarshal(msg, &analyticsMsg); err != nil {
			logger.Error("failed to unmarshal message", zap.Error(err))
			continue
		}

		output <- analyticsMsg
	}
}

func (p *Port) Close() error {
	return p.consumer.Close()
}
