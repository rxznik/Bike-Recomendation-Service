package parsing

import (
	"context"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type ParsingPort interface {
	Accept(msg models.AnalyticsMessage) models.ParsingInput
}

type ParsingExternalAdapter interface {
	SendToShared(outputMsg models.ParsingOutput, shared *models.RecommendationsMessage)
}

type ParsingParseAdapter interface {
	GetRelevantProduct(ctx context.Context, detail string) (*models.ParsingParseResponse, error)
}

type ParsingApp struct {
	externalAdapter ParsingExternalAdapter
	parseAdapter    ParsingParseAdapter
	port            ParsingPort
	log             *zap.Logger
}

func New(
	logger *zap.Logger,
	port ParsingPort,
	externalAdapter ParsingExternalAdapter,
	parseAdapter ParsingParseAdapter,
) *ParsingApp {
	logger = logger.With(zap.String("service", "parsing"))
	return &ParsingApp{
		externalAdapter: externalAdapter,
		parseAdapter:    parseAdapter,
		port:            port,
		log:             logger,
	}
}

func (p *ParsingApp) ProcessMessage(ctx context.Context, msg models.AnalyticsMessage, shared *models.RecommendationsMessage) error {
	ctx, span := otel.Tracer("parsing").Start(ctx, "ProcessMessage")
	defer span.End()
	parsingInput := p.port.Accept(msg)
	p.log.Info("parsing input", zap.String("detail", parsingInput.Detail))

	adapterCtx, adapterSpan := otel.Tracer("parsing").Start(ctx, "GetRelevantProduct")
	relevantProductFromMarket, err := p.parseAdapter.GetRelevantProduct(adapterCtx, parsingInput.Detail)
	if err != nil {
		adapterSpan.RecordError(err)
		adapterSpan.SetStatus(codes.Error, err.Error())
		return err
	}
	adapterSpan.End()
	if relevantProductFromMarket == nil {
		p.log.Info("no relevant product found", zap.String("detail", parsingInput.Detail))
		return nil
	}
	parsingOutput := models.ParsingOutput{Market: relevantProductFromMarket.URL}
	p.externalAdapter.SendToShared(parsingOutput, shared)
	return nil
}
