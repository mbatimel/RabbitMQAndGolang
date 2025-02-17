package service

import (
	"context"
	"fmt"

	externalApi "github.com/mbatimel/RabbitMQAndGolang/limits/internal/interfaces"
	"github.com/mbatimel/RabbitMQAndGolang/limits/internal/models"
	"github.com/rs/zerolog"
)

type UnitOfWork interface {
	Rollback(context.Context) error
	Commit(context.Context) error
}
type Storage interface {
	GetUnitOfWork(context.Context, bool) (UnitOfWork, error)
	AddLimits(ctx context.Context, ouw UnitOfWork, count int) (err error)
}
type subscriptionService struct {
	logger  zerolog.Logger
	storage Storage
}

func (s *subscriptionService) AddLimits(ctx context.Context) (err error) {
	uow, err := s.storage.GetUnitOfWork(ctx, models.MASTER)
	if err != nil {
		return fmt.Errorf("could not obtain unit of work: %w", err)
	}
	defer func() {
		_ = uow.Rollback(ctx)
	}()
// TODO: дописать логику выдачи лимитов в зависимости от подписки
	return s.storage.AddLimits(ctx, uow, count)

}
func Newservice(logger zerolog.Logger, storage Storage) externalApi.Subscription {
	return &subscriptionService{
		storage: storage,
		logger:  logger,
	}
}
