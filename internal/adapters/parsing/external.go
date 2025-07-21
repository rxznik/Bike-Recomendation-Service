package parsing

import (
	"sync"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.uber.org/zap"
)

type ParsingExternalAdapter struct {
	mu  sync.Mutex
	log *zap.Logger
}

func NewExternal(logger *zap.Logger) *ParsingExternalAdapter {
	logger = logger.With(zap.String("adapter", "parsing"))
	return &ParsingExternalAdapter{log: logger}
}

func (p *ParsingExternalAdapter) SendToShared(msg models.ParsingOutput, shared *models.RecommendationsMessage) {
	p.mu.Lock()
	shared.Recomendations.Market = msg.Market
	p.mu.Unlock()
}
