package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5"
)

//go:embed sql/insert_active_sub.sql
var insertActiveSubSql string

func (s *storage) activateSubscription(ctx context.Context, tx pgx.Tx, tariffId int, price int) (err error) {
	supplierID := rand.Intn(100)
	_, err = tx.Exec(ctx, insertActiveSubSql, supplierID, tariffId, price)
	if err != nil {
		return fmt.Errorf("postgresql insert new active sub error: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
