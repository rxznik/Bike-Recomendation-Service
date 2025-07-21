package setup

import (
	"net/http"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/adapters"
	geoadapter "github.com/devprod-tech/webike_recomendations-Vitalya/internal/adapters/geo"
	parsingadapter "github.com/devprod-tech/webike_recomendations-Vitalya/internal/adapters/parsing"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/app"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/config"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/ports"
	geoport "github.com/devprod-tech/webike_recomendations-Vitalya/internal/ports/geo"
	parsingport "github.com/devprod-tech/webike_recomendations-Vitalya/internal/ports/parsing"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/queue/rabbitmq"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/services/geo"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/services/parsing"
	"go.uber.org/zap"
)

type StatusCommitter interface {
	CommitMessageStatus(queue string, status string)
	CommitWorkStatus(status string)
}

func MustSetupRabbitMQApplication(logger *zap.Logger, cfg *config.Config, statusCommitter StatusCommitter) *app.App {
	port := ports.NewRabbitMQ(logger, &ports.RabbitMQPortOptions{
		QueueOptions: &rabbitmq.QueueOptions{
			Name:    cfg.Consumer.QueueName,
			Durable: true,
		},
		RabbitMQDSN:            cfg.BrokerDSN,
		MessageStatusCommitter: statusCommitter,
	})
	adapter := adapters.NewRabbitMQ(logger, &adapters.RabbitMQAdapterOptions{
		QueueOptions: &rabbitmq.QueueOptions{
			Name:    cfg.Consumer.QueueName,
			Durable: true,
		},
		RabbitMQDSN:            cfg.BrokerDSN,
		MessageStatusCommitter: statusCommitter,
	})
	services := MustSetupServices(logger, cfg, statusCommitter)
	return app.New(logger, port, adapter, services...)
}

func MustSetupServices(logger *zap.Logger, cfg *config.Config, adaptersWorkStatusCommitter adapters.WorkStatusCommitter) []app.Service {
	client := &http.Client{Timeout: cfg.HttpClientTimeout}

	geoPort := geoport.New(logger)
	parsingPort := parsingport.New(logger)

	geoExternalAdapter := geoadapter.NewExternal(logger)
	geoGoogleAdapter := geoadapter.NewGoogle(logger, client, cfg.GoogleAPIKey)

	geoGoogleAdapter.SetupWorkStatusCommitter(adaptersWorkStatusCommitter)

	parsingExternalAdapter := parsingadapter.NewExternal(logger)
	parsingParseAdapter := parsingadapter.NewParse(logger, client)

	parsingParseAdapter.SetupWorkStatusCommitter(adaptersWorkStatusCommitter)

	geoApp := geo.New(logger, geoPort, geoExternalAdapter, geoGoogleAdapter)
	parsingApp := parsing.New(logger, parsingPort, parsingExternalAdapter, parsingParseAdapter)
	return []app.Service{geoApp, parsingApp}
}
