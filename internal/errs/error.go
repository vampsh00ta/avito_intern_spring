package errs

import (
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

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf(DublicateErr)
		}
	}
	if err == pgx.ErrNoRows {
		return fmt.Errorf(NoRowsInResultErr)

	}
	return err
}
