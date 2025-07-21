package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	StatusMessageProduced      = "produced"
	StatusMessageProduceFailed = "produce_failed"
	StatusMessageConsumed      = "consumed"
	StatusMessageConsumeFailed = "consume_failed"
	StatusMarketParsed         = "market_parsing"
	StatusMarketError          = "market_error"
	StatusNearestTOError       = "to_error"
	StatusNearestTOFound       = "to_found"
)

var (
	messageConsumedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "recomendation_app_messages_consumed_total",
		Help: "Total number of processed events",
	}, []string{"queue", "status"})

	messageProducedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "recomendation_app_messages_produced_total",
		Help: "Total number of produced events",
	}, []string{"queue", "status"})

	marketParsedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "recomendation_app_market_parsed_total",
		Help: "Total number of market parsed events",
	}, []string{"status"})

	nearestToFoundCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "recomendation_app_nearest_to_found_total",
		Help: "Total number of nearest TO found events",
	}, []string{"status"})
)

type Prometheus struct {
	log                    *zap.Logger
	srv                    *http.Server
	messageConsumedCounter *prometheus.CounterVec
	messageProducedCounter *prometheus.CounterVec
	marketParsedCounter    *prometheus.CounterVec
	nearestToFoundCounter  *prometheus.CounterVec
}

func NewPrometheus(logger *zap.Logger) *Prometheus {
	logger = logger.With(zap.String("observability", "prometheus"))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &Prometheus{
		log:                    logger,
		srv:                    &http.Server{Addr: ":2112", Handler: mux},
		messageConsumedCounter: messageConsumedCounter,
		messageProducedCounter: messageProducedCounter,
		marketParsedCounter:    marketParsedCounter,
		nearestToFoundCounter:  nearestToFoundCounter,
	}
}

func (p *Prometheus) CommitMessageStatus(queue, status string) {
	switch status {
	case StatusMessageConsumed, StatusMessageConsumeFailed:
		p.messageConsumedCounter.WithLabelValues(queue, status).Inc()
	case StatusMessageProduced, StatusMessageProduceFailed:
		p.messageProducedCounter.WithLabelValues(queue, status).Inc()
	}
}

func (p *Prometheus) CommitWorkStatus(status string) {
	switch status {
	case StatusMarketParsed, StatusMarketError:
		p.marketParsedCounter.WithLabelValues(status).Inc()
	case StatusNearestTOFound, StatusNearestTOError:
		p.nearestToFoundCounter.WithLabelValues(status).Inc()
	}
}

func (p *Prometheus) Run() error {
	p.log.Info("starting prometheus server")
	if err := p.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		p.log.Fatal("failed to start prometheus server", zap.Error(err))
		return err
	}
	return nil
}

func (p *Prometheus) Stop() error {
	return p.srv.Close()
}
