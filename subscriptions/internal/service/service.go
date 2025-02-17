package service

import (
	"context"
	"fmt"

	externalApi "github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/interfaces"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models"
	"github.com/rs/zerolog"
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
	logger  zerolog.Logger
	storage Storage
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
	return s.storage.ActiveSubscription(ctx, uow, limitId, price)

}
func Newservice(logger zerolog.Logger, storage Storage) externalApi.Subscription {
	return &subscriptionService{
		storage: storage,
		logger:  logger,
	}
}
