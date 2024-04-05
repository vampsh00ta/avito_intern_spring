package repository

import (
	"avito_intern/internal/repository/psql"
	"avito_intern/internal/repository/redis"
)

type Repository interface {
	psql.Repository
	redis.Repository
}

type cache struct {
}
type repository struct {
	//db
	cache
}

func New() Repository {
	return repository{}
}
