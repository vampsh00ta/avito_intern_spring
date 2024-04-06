package psql

import (
	"avito_intern/pkg/client"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Tx interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	getDb(ctx context.Context) (client.Client, error)
}

func (db Pg) Rollback(ctx context.Context) error {
	tx, err := db.getTx(ctx)
	if err != nil {
		return err
	}
	if err := tx.Rollback(ctx); err != nil {
		return err
	}
	return nil
}
func (db Pg) Commit(ctx context.Context) error {
	tx, err := db.getTx(ctx)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
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
		return nil, fmt.Errorf("no tx")
	}
	return txModel, nil
}
func (db Pg) getDb(ctx context.Context) (client.Client, error) {
	tx := ctx.Value("tx")
	txModel, ok := tx.(client.Client)
	if !ok {
		txModel = db.client
	}
	return txModel, nil
}
