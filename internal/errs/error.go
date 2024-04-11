package errs

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	IncorrectJSONErr  = "incorrect json"
	DublicateErr      = "some of input data already exists"
	NilIDErr          = "nil ID"
	WrongIDErr        = "wrong id"
	IncorrectTokenErr = "incorrect token"
	AuthErr           = "auth error"
	InvalidTokenErr   = "invalid token"
	NotAdminErr       = "you are not admin"
	NoUserSuchUserErr = "no such user"
	NotLoggedErr      = "you are not logged"
	WrongRoleErr      = "wrong role"
	NoRowsInResultErr = "no data affected"
)

func Handle(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf(NoRowsInResultErr)

	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			return fmt.Errorf(DublicateErr)

		}
	}

	return err
}
