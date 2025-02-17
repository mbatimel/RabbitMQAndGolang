package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

const RETRY_INSERT_LIMIT = 5

//go:embed sql/insert_active_sub.sql
var insertActiveSubSql string

func (s *storage) activateSubscription(ctx context.Context, tx pgx.Tx, tariffId int, price int) (err error) {
	_, err := tx.Exec(ctx, insertActiveSubSql, tariffId, price)
	if err != nil {
		return fmt.Errorf("postgresql insert new active sub error: %w", err)
	}

	return nil
}
