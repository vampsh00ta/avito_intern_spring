package errs

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"

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
	NoReferenceErr    = "no such tag/feature"

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

	case errors.Is(err, jwt.ErrSignatureInvalid), errors.Is(err, jwt.ErrTokenMalformed):
		return fmt.Errorf(AuthErr)

	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf(DublicateErr)
		case "23503":
			return fmt.Errorf(NoReferenceErr)
		}

	}

	return err
}
