package app

import (
	"context"
	"reflect"
	"sync"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/observability"
	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Port interface {
	Open() (<-chan models.AnalyticsMessage, error)
	Close() error
}

type Adapter interface {
	Produce(ctx context.Context, message models.RecommendationsMessage) error
	Close() error
}

type Service interface {
	ProcessMessage(ctx context.Context, msg models.AnalyticsMessage, output *models.RecommendationsMessage) error
}

type App struct {
	wg       sync.WaitGroup
	port     Port
	adapter  Adapter
	services []Service
	log      *zap.Logger
}

func New(logger *zap.Logger, port Port, adapter Adapter, services ...Service) *App {
	return &App{
		adapter:  adapter,
		port:     port,
		services: services,
		log:      logger,
	}
}

func (a *App) Run(ctx context.Context) error {

	ctx, span := otel.Tracer("app").Start(ctx, "Run")
	defer span.End()

	appMsgChan, err := a.port.Open()
	if err != nil {
		a.log.Fatal("failed to open port", zap.Error(err))
		observability.CaptureExceptionSentry(err)
		return err
	}
	for msg := range appMsgChan {
		a.wg.Add(1)
		go func(ctx context.Context) {
			prcessMsgAndProduceCtx, prcessMsgAndProduceSpan := otel.Tracer("app").Start(ctx, "ProcessMessageAndProduce")
			defer func() {
				a.wg.Done()
				prcessMsgAndProduceSpan.End()
			}()
			a.processMsgAndProduce(prcessMsgAndProduceCtx, msg)
		}(ctx)
	}

	a.wg.Wait()

	return nil
}

func (a *App) processMsgAndProduce(ctx context.Context, msg models.AnalyticsMessage) {
	var wg sync.WaitGroup

	output := models.RecommendationsMessage{
		AnalyticsMessage: msg,
	}

	for _, service := range a.services {
		wg.Add(1)
		go func(localHub *sentry.Hub, service Service) {
			svcCtx, svcSpan := otel.Tracer("app/service").Start(ctx, reflect.TypeOf(service).Name())
			defer svcSpan.End()
			defer wg.Done()
			if err := service.ProcessMessage(svcCtx, msg, &output); err != nil {
				localHub.CaptureException(err)
				svcSpan.RecordError(err)
				svcSpan.SetStatus(codes.Error, err.Error())
				a.log.Error("failed to process message", zap.Error(err))
			}
		}(observability.NewLocalHubSentry(), service)
	}

	wg.Wait()
	produceCtx, produceSpan := otel.Tracer("app").Start(ctx, "Produce")
	defer produceSpan.End()
	if err := a.adapter.Produce(produceCtx, output); err != nil {
		observability.CaptureExceptionSentry(err)
		produceSpan.RecordError(err)
		produceSpan.SetStatus(codes.Error, err.Error())
		a.log.Fatal("failed to produce message", zap.Error(err))
	}
}

func (a *App) Stop() error {
	if err := a.port.Close(); err != nil {
		observability.CaptureExceptionSentry(err)
		a.log.Fatal("failed to close port", zap.Error(err))
		return err
	}
	if err := a.adapter.Close(); err != nil {
		observability.CaptureExceptionSentry(err)
		a.log.Fatal("failed to close adapter", zap.Error(err))
		return err
	}
	return nil
}
