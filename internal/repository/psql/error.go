package psql

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func handleEror(err error) error {

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			err = fmt.Errorf("er")
		}
	}
	return err
}
