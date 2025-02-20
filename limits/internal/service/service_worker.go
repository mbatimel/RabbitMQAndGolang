package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
	AddLimits(ctx context.Context, ouw UnitOfWork, supplierIDint, count int, limit_id int, describe string) (err error)
}
type LimitsWorker struct {
	logger       zerolog.Logger
	storage      Storage
	ctx          context.Context
	metric       *metrics.Metrics
	cancelFunc   context.CancelFunc
	rabbitMQConn *amqp.Connection
}

func (s *LimitsWorker) AddLimits(ctx context.Context, supplierID int, count int, limit_id int, describe string) (err error) {
	uow, err := s.storage.GetUnitOfWork(ctx, models.MASTER)
	if err != nil {
		return fmt.Errorf("could not obtain unit of work: %w", err)
	}
	defer func() {
		_ = uow.Rollback(ctx)
	}()

	return s.storage.AddLimits(ctx, uow, supplierID, count, limit_id, describe)

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

	s.logger.With().Str("msgID: %s \nReceived message from RabbitMQ queue: %s", string(msg.MessageId)).Logger()
	parts := strings.Split(string(msg.Body), ":")
	var supplierID int
	var limit_id int
	var count int
	var describe string
	supplierID, err := strconv.Atoi(parts[0])
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to parse supplierID")
	}
	limit_id, err = strconv.Atoi(parts[1])
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to parse limit id")
	}
	if limit_id == 1 {
		count = 100
		describe = "топовый лимит"
	}
	if err := s.AddLimits(s.ctx, supplierID, limit_id, count, describe); err != nil {
		s.logger.Error().Err(err).Msg("Failed to parse limit id")
	}
	// Подтверждаем получение сообщения
	err = msg.Ack(false)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to acknowledge message")
	}

}

func (s *LimitsWorker) StopWorker() {
	s.cancelFunc()
	_ = s.rabbitMQConn.Close()
}
