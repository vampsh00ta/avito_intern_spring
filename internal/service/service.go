package service

import (
	//"avito_intern/config"
	psqlrepo "avito_intern/internal/repository/psql"
	redisrepo "avito_intern/internal/repository/redis"
)

type Service interface {
	Banner
}
type service struct {
	db    psqlrepo.Repository
	cache redisrepo.Repository
}

func New(psqlrepo psqlrepo.Repository, redisrepo redisrepo.Repository) Service {
	return &service{
		psqlrepo,
		redisrepo,
	}
}
