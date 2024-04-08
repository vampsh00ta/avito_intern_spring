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
}

func (db Pg) Rollback(ctx context.Context) error {
	txCtx := ctx.Value("tx")
	tx, ok := txCtx.(pgx.Tx)
	if !ok {
		return fmt.Errorf("no tx")
	}

	if err := tx.Rollback(ctx); err != nil {
		return err
	}
	return nil
}
func (db Pg) Commit(ctx context.Context) error {
	txCtx := ctx.Value("tx")
	tx, ok := txCtx.(pgx.Tx)
	if !ok {
		return fmt.Errorf("no tx")
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

	ctx = context.WithValue(ctx, "tx", tx)
	ctx = context.WithValue(ctx, "client", client.Client(tx))

	return ctx, nil
}

func (db Pg) getDb(ctx context.Context) (client.Client, error) {
	tx := ctx.Value("client")
	txModel, ok := tx.(client.Client)
	if !ok {
		txModel = db.client
	}
	return txModel, nil
}
