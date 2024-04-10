package service

import (
	"avito_intern/config"
	//"avito_intern/config".
	psqlrepo "avito_intern/internal/repository/psql"
	redisrepo "avito_intern/internal/repository/redis"
)

type Service interface {
	Banner
	BannerHistory
	Auth
}
type service struct {
	db    psqlrepo.Repository
	cache redisrepo.Repository
	cfg   *config.Config
}

func New(psqlrepo psqlrepo.Repository, redisrepo redisrepo.Repository, cfg *config.Config) Service {
	srvc := &service{
		psqlrepo,
		redisrepo,
		cfg,
	}
	srvc.bannerHistoryCleaner(3)
	return srvc
}
