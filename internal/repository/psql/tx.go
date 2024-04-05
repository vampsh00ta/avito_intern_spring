package psql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type Tx interface {
	Begin(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	getTx(ctx context.Context) *pgx.Tx
}

func (db Pg) Rollback(ctx context.Context) error {
	tx, err := db.getTx(ctx)
	if err != nil {
		return err
	}
	tx.Rollback(ctx)
	return nil
}
func (db Pg) Commit(ctx context.Context) error {
	tx, err := db.getTx(ctx)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}
func (db Pg) Begin(ctx context.Context) (context.Context, error) {
	tx, err := db.client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	txCtx := context.WithValue(ctx, "tx", tx)
	return txCtx, nil
}
func (db Pg) getTx(ctx context.Context) (pgx.Tx, error) {
	tx := ctx.Value("tx")
	txModel, ok := tx.(pgx.Tx)
	if !ok {
		return nil, errors.New("no tx")
	}
	return txModel, nil
}
