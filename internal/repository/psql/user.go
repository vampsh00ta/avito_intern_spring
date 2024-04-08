package psql

import (
	"avito_intern/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type User interface {
	GetUserById(ctx context.Context, id int) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}

// select * from customer where id = $1 limit  1
func (db Pg) GetUserById(ctx context.Context, id int) (models.User, error) {
	q := `select * from customer where id = $1 limit  1`
	client, err := db.getDb(ctx)
	if err != nil {
		return models.User{}, err
	}

	row, err := client.Query(ctx, q, id)
	if err != nil {
		return models.User{}, err
	}
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// select * from customer where username = $1 limit  1
func (db Pg) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	q := `select * from customer where username = $1 limit  1`
	client, err := db.getDb(ctx)
	if err != nil {
		return models.User{}, err
	}

	row, err := client.Query(ctx, q, username)
	if err != nil {
		return models.User{}, err
	}
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
