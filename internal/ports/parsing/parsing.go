package parsing

import (
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.uber.org/zap"
)

type ParsingPort struct {
	log *zap.Logger
}

func New(logger *zap.Logger) *ParsingPort {
	logger = logger.With(zap.String("port", "parsing"))
	return &ParsingPort{log: logger}
}

func (p *ParsingPort) Accept(msg models.AnalyticsMessage) models.ParsingInput {
	return models.ParsingInput{Detail: msg.Detail}
}
