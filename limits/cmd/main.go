package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/config"
	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/metrics"
	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/service"
	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/storage/postgres"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

const serviceName = "limits"

func main() {
	log.Logger = config.Values().Logger().With().Str("serviceName", serviceName).Logger()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	// Подключение к Postgres
	postgresStorage, err := postgres.New(log.Logger, config.Values().Postgres)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	storage := postgres.NewStorage(postgresStorage)

	// Подключение к RabbitMQ
	rabbitMQConn, err := amqp.Dial(config.Values().RabbitMQ.URL)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to connect to RabbitMQ")
	}

	metricsCollector := metrics.CreateMetrics(serviceName, serviceName)
	ctx := context.Background()

	// Передаем RabbitMQ соединение в сервис
	svc := service.Newservice(ctx, log.Logger, storage, metricsCollector, rabbitMQConn)

	go func() {
		log.Info().Msg("started consumer service")
		svc.StartWorker(config.Values().RabbitMQ.Queue)
	}()

	<-shutdown
	svc.StopWorker()
}
