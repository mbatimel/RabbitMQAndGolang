package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const RETRY_INSERT_LIMIT = 5

//go:embed sql/insert_limits_for_sub.sql
var insertlimitsSql string

func (s *storage) addLimits(ctx context.Context, tx pgx.Tx, count int, limit_id int, describe string) (err error) {
	_, err = tx.Exec(ctx, insertlimitsSql, count, limit_id, describe)
	if err != nil {
		return fmt.Errorf("postgresql insert new active sub error: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
