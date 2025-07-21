package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/config"
	mockproducer "github.com/devprod-tech/webike_recomendations-Vitalya/internal/mock-producer"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/queue/rabbitmq"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/setup"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const numProducers = 10

var (
	logger        *zap.Logger
	cfg           *config.Config
	mockProducers = make([]*mockproducer.MockProducer, 0, numProducers)
)

func init() {
	godotenv.Load()
	cfg = config.MustLoad()
	logger = setup.MustSetupLogger(cfg.Env)
	for range numProducers {
		producer := rabbitmq.NewConnection(logger, cfg.BrokerDSN)
		prodOpts := rabbitmq.QueueOptions{Name: cfg.Consumer.QueueName, Durable: true}
		if err := producer.DeclareQueue(prodOpts); err != nil {
			logger.Fatal("failed to declare a queue", zap.Error(err))
		}
		mockProducer := mockproducer.New(logger, producer)
		mockProducers = append(mockProducers, mockProducer)
	}
}

func main() {
	defer logger.Sync()

	logger.Info("starting mock producers", zap.Int("num_producers", numProducers))
	for i, mockProducer := range mockProducers {
		go func() {
			logger.Info("mock producer started", zap.Int("id", i))
			if err := mockProducer.Run(); err != nil {
				logger.Fatal("failed to run mock producer", zap.Error(err))
			}
		}()
		time.Sleep(1 * time.Second)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("shutting down mock producers")
	for _, producer := range mockProducers {
		if err := producer.Stop(); err != nil {
			logger.Fatal("failed to stop mock producer", zap.Error(err))
		}
	}
	logger.Info("mock producers shutdown")
}
