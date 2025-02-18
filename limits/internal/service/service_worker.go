package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/metrics"
	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/models"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type UnitOfWork interface {
	Rollback(context.Context) error
	Commit(context.Context) error
}
type Storage interface {
	GetUnitOfWork(context.Context, bool) (UnitOfWork, error)
	AddLimits(ctx context.Context, ouw UnitOfWork, count int, limit_id int, describe string) (err error)
}
type LimitsWorker struct {
	logger       zerolog.Logger
	storage      Storage
	ctx          context.Context
	metric       *metrics.Metrics
	cancelFunc   context.CancelFunc
	rabbitMQConn *amqp.Connection
}

func (s *LimitsWorker) AddLimits(ctx context.Context, count int, limit_id int, describe string) (err error) {
	uow, err := s.storage.GetUnitOfWork(ctx, models.MASTER)
	if err != nil {
		return fmt.Errorf("could not obtain unit of work: %w", err)
	}
	defer func() {
		_ = uow.Rollback(ctx)
	}()

	return s.storage.AddLimits(ctx, uow, count, limit_id, describe)

}
func Newservice(ctx context.Context, logger zerolog.Logger, storage Storage, metricsCollector *metrics.Metrics, rabbitMQConn *amqp.Connection) *LimitsWorker {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &LimitsWorker{
		ctx:          ctx,
		storage:      storage,
		logger:       logger,
		metric:       metricsCollector,
		cancelFunc:   cancelFunc,
		rabbitMQConn: rabbitMQConn,
	}

}
func (s *LimitsWorker) StartWorker(queueName string) {
	ch, err := s.rabbitMQConn.Channel()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to open a channel")
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // имя очереди
		"",        // consumer tag
		false,     // auto-ack (нам нужно вручную подтверждать)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to register consumer")
		return
	}

	s.logger.Info().Msg("RabbitMQ consumer started")

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info().Msg("RabbitMQ consumer stopped")
			return

		case msg, ok := <-msgs:
			if !ok {
				s.logger.Warn().Msg("RabbitMQ channel closed")
				return
			}
			s.HandleMessage(msg)

		case <-time.After(5 * time.Second):
			s.logger.Info().Msg("No messages in RabbitMQ queue")
		}
	}
}

func (s *LimitsWorker) HandleMessage(msg amqp.Delivery) {

	logger := s.logger.With().Str("msgID", string(msg.MessageId)).Logger()

	// Логируем сообщение
	logger.Info().Msgf("Received message from RabbitMQ queue: %s", string(msg.Body))

	// Подтверждаем получение сообщения
	err := msg.Ack(false)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to acknowledge message")
	}

}

func (s *LimitsWorker) StopWorker() {
	s.cancelFunc()
	_ = s.rabbitMQConn.Close()
}
