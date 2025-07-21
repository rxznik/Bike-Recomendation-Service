package geo

import (
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.uber.org/zap"
)

type GeoPort struct {
	log *zap.Logger
}

func New(logger *zap.Logger) *GeoPort {
	logger = logger.With(zap.String("port", "geo"))
	return &GeoPort{log: logger}
}

func (g *GeoPort) Accept(msg models.AnalyticsMessage) models.GeoInput {
	return models.GeoInput{
		Latitude:  msg.Payload.Location.Latitude,
		Longitude: msg.Payload.Location.Longitude,
	}
}
