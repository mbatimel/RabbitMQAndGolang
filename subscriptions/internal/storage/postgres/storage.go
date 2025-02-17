package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/service"
)

var _ service.Storage = &storage{}

type storage struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool
}

func (s *storage) GetUnitOfWork(ctx context.Context, isMaster bool) (service.UnitOfWork, error) {
	if isMaster {
		return s.master.Begin(ctx)
	}
	return s.replica.Begin(ctx)
}
func (s *storage) unpackUnitOfWork(unitOfWork service.UnitOfWork) (pgx.Tx, error) {
	tx, ok := unitOfWork.(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("unit of work must be a pgx transaction")
	}
	return tx, nil
}
func (s *storage) ActiveSubscription(ctx context.Context, ouw service.UnitOfWork, limitId int, price int) (err error) {
	tx, err := s.unpackUnitOfWork(ouw)
	if err != nil {
		return fmt.Errorf("could not unpack uow: %w", err)
	}
	return s.activateSubscription(ctx, tx, limitId, price)
}
func NewStorage(conn ConnectManager) *storage {
	master, replica := conn.GetLimitsConn()
	return &storage{
		master:  master,
		replica: replica,
	}
}
