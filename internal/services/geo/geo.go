package geo

import (
	"context"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type GeoPort interface {
	Accept(msg models.AnalyticsMessage) models.GeoInput
}

type GeoExternalAdapter interface {
	SendToShared(outputMsg models.GeoOutput, shared *models.RecommendationsMessage)
}

type GeoAPIAdapter interface {
	GetNearestTO(ctx context.Context, longitude float64, latitude float64) (*models.GeoAPIResponse, error)
}

type GeoApp struct {
	externalAdapter GeoExternalAdapter
	apiAdapter      GeoAPIAdapter
	port            GeoPort
	log             *zap.Logger
}

func New(
	logger *zap.Logger,
	port GeoPort,
	externalAdapter GeoExternalAdapter,
	apiAdapter GeoAPIAdapter,
) *GeoApp {
	logger = logger.With(zap.String("service", "geo"))
	return &GeoApp{
		port:            port,
		externalAdapter: externalAdapter,
		apiAdapter:      apiAdapter,
		log:             logger,
	}
}

func (g *GeoApp) ProcessMessage(ctx context.Context, msg models.AnalyticsMessage, shared *models.RecommendationsMessage) error {
	ctx, span := otel.Tracer("geo").Start(ctx, "ProcessMessage")
	defer span.End()
	geoInput := g.port.Accept(msg)
	g.log.Info("geo input", zap.Float64("latitude", geoInput.Latitude), zap.Float64("longitude", geoInput.Longitude))
	adapterCtx, adapterSpan := otel.Tracer("geo").Start(ctx, "GetNearestTO")
	data, err := g.apiAdapter.GetNearestTO(adapterCtx, geoInput.Longitude, geoInput.Latitude)
	if err != nil {
		adapterSpan.RecordError(err)
		adapterSpan.SetStatus(codes.Error, err.Error())
		return err
	}
	adapterSpan.End()
	if data == nil {
		g.log.Info("no nearest TO found", zap.Float64("latitude", geoInput.Latitude), zap.Float64("longitude", geoInput.Longitude))
		return nil
	}
	nearestTO := data.DisplayName
	g.log.Info("finding nearest TO", zap.String("TO", nearestTO))
	geoOutput := models.GeoOutput{Nearest_TO: nearestTO}
	g.externalAdapter.SendToShared(geoOutput, shared)
	return nil
}
