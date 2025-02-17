package main

// import (
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"time"

// 	"syscall"

// 	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/config"
// 	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/service"
// 	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/storage/postgres"
// 	transportHttp "github.com/mbatimel/RabbitMQAndGolang/limits/internal/transport/http"
// 	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/transport/jsonRPC/externalapi"
// 	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/transport/jsonRPC/middlewares"
// 	"github.com/rs/zerolog/log"
// 	"github.com/valyala/fasthttp"
// )

// const serviceName = "limits"

// func main() {
// 	log.Logger = config.Values().Logger().With().Str("serviceName", serviceName).Logger()
// 	shutdown := make(chan os.Signal, 1)
// 	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

// 	postgresStorage, err := postgres.New(log.Logger,config.Values().Postgres)
// 	if err != nil {
// 		log.Logger.Fatal().Err(err).Msg("failed to connect to postgres")
// 	}
// 	storage := postgres.NewStorage(postgresStorage)
// 	svc := service.Newservice(log.Logger, storage)

// 	services := []externalapi.Option{
// 		externalapi.Use(middlewares.Recover),
// 		externalapi.Subscription(externalapi.NewSubscription(svc)),
// 	}
// 	app := externalapi.New(log.Logger, services...).WithLog().WithMetrics()
// 	server := &fasthttp.Server{
// 		Handler:            app.Fiber().Handler(),
// 		MaxRequestBodySize: config.Values().MaxRequestBodySize,
// 		ReadBufferSize:     config.Values().MaxRequestHeaderSize,
// 		ReadTimeout:        time.Duration(config.Values().ReadTimeout) * time.Second,
// 	}

// 	healthServer := transportHttp.NewHealthServer()

// 	wg := &sync.WaitGroup{}

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		app.ServeMetrics(log.Logger, config.Values().MetricsPath, config.Values().MetricsBind)
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		serveErr := server.ListenAndServe(config.Values().ServiceBind)
// 		if serveErr != nil {
// 			log.Fatal().Err(serveErr).Msg("failed to listen and serve pay-api-internal server")
// 		} else {
// 			log.Error().Msg("external api pay-api-internal server stopped with no error")
// 		}
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		healthErr := healthServer.Start(config.Values().HealthBind)
// 		if healthErr != nil {
// 			log.Error().Err(healthErr).Msg("failed to start health server")
// 		} else {
// 			log.Error().Msg("health server stopped with no error")
// 		}
// 	}()

// 	<-shutdown
// 	err = healthServer.Stop()
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to stop health server")
// 	}

// 	err = server.Shutdown()
// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to shutdown server")
// 	}

// 	wg.Wait()
// }
