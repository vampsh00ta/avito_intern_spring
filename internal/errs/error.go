package errs

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	IncorrectJSONErr  = errors.New("incorrect json")
	DublicateErr      = errors.New("some of input data already exists")
	NilIDErr          = errors.New("nil ID")
	WrongIDErr        = errors.New("wrong id")
	IncorrectTokenErr = errors.New("incorrect token")
	AuthErr           = errors.New("auth error")
	InvalidTokenErr   = errors.New("invalid token")
	NoReferenceErr    = errors.New("no such tag/feature")
	UnknownErr        = errors.New("unknown error")
	ValidationError   = errors.New("incorrect input data")

	NotAdminErr       = errors.New("you are not admin")
	NoUserSuchUserErr = errors.New("no such user")
	NotLoggedErr      = errors.New("you are not logged")
	WrongRoleErr      = errors.New("wrong role")
	NoRowsInResultErr = errors.New("no such data")
)

func Handle(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return NoRowsInResultErr

	case errors.Is(err, jwt.ErrSignatureInvalid), errors.Is(err, jwt.ErrTokenMalformed):
		return AuthErr

	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			return DublicateErr
		case "23503":
			return NoReferenceErr
		default:
			return UnknownErr
		}
	}
	if _, ok := err.(validator.ValidationErrors); ok {
		return ValidationError
	}
	if _, ok := err.(redis.Error); ok {
		return UnknownErr
	}

	return err
}
