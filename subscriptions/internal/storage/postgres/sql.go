package postgres

import (
	"context"
	_ "embed"
	"fmt"


	"github.com/jackc/pgx/v5"

)

const RETRY_INSERT_LIMIT = 5

//go:embed sql/insert_active_sub.sql
var insertActiveSubSql string

func (s *storage) activateSubscription(ctx context.Context, tx pgx.Tx, tariffId int, price int) (err error) {
	_, err = tx.Exec(ctx, insertActiveSubSql, tariffId, price)
	if err != nil {
		return fmt.Errorf("postgresql insert new active sub error: %w", err)
	}

	return nil
}
