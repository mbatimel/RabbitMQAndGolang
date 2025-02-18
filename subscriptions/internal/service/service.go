package service

import (
	"context"
	"fmt"

	externalApi "github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/interfaces"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type UnitOfWork interface {
	Rollback(context.Context) error
	Commit(context.Context) error
}
type Storage interface {
	GetUnitOfWork(context.Context, bool) (UnitOfWork, error)
	ActiveSubscription(ctx context.Context, ouw UnitOfWork, limitId int, price int) (err error)
}
type subscriptionService struct {
	logger   zerolog.Logger
	storage  Storage
	rabbitMQ *amqp.Connection
}

func (s *subscriptionService) ActiveSubscription(ctx context.Context, limitId int, price int) (err error) {
	uow, err := s.storage.GetUnitOfWork(ctx, models.MASTER)
	if err != nil {
		return fmt.Errorf("could not obtain unit of work: %w", err)
	}
	defer func() {
		_ = uow.Rollback(ctx)
	}()
	if limitId == 0 {
		return fmt.Errorf("Why are you don't have a limitID ((((")
	}
	if err := s.storage.ActiveSubscription(ctx, uow, limitId, price); err != nil {
		return fmt.Errorf("error activate subscription: %w", err)
	}

	ch, err := s.rabbitMQ.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"Subscription_limits",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare: %w", err)
	}

	err = ch.Publish(
		"",
		"Subscription_limits",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")
	return nil

}
func Newservice(logger zerolog.Logger, storage Storage, rabbitMQ *amqp.Connection) externalApi.Subscription {
	return &subscriptionService{
		storage:  storage,
		logger:   logger,
		rabbitMQ: rabbitMQ,
	}
}
