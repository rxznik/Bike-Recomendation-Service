package geo

import (
	"sync"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.uber.org/zap"
)

type GeoExternalAdapter struct {
	mu  sync.Mutex
	log *zap.Logger
}

func NewExternal(logger *zap.Logger) *GeoExternalAdapter {
	logger = logger.With(zap.String("adapter", "geo"))
	return &GeoExternalAdapter{log: logger}
}

func (g *GeoExternalAdapter) SendToShared(msg models.GeoOutput, shared *models.RecommendationsMessage) {
	g.mu.Lock()
	shared.Recomendations.Nearest_TO = msg.Nearest_TO
	g.mu.Unlock()
}
