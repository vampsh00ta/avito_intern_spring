package service

import (
	"avito_intern/config"

	psqlrepo "avito_intern/internal/repository/psql"
	redisrepo "avito_intern/internal/repository/redis"

	"go.uber.org/zap"
	//"avito_intern/config".
)

type Service interface {
	Banner
	BannerHistory
	Auth
}
type service struct {
	db     psqlrepo.Repository
	cache  redisrepo.Repository
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func New(psqlrepo psqlrepo.Repository,
	redisrepo redisrepo.Repository,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) Service {
	srvc := &service{
		psqlrepo,
		redisrepo,
		cfg,
		logger,
	}
	doneBannerHistoryCleaner := make(chan bool)
	msgs := make(chan error)

	srvc.BannerHistoryCleaner(msgs, doneBannerHistoryCleaner, 3)
	go func(msgs <-chan error) {
		for err := range msgs {
			if err != nil {
				srvc.logger.Error(err)
			} else {
				srvc.logger.Info("cleared history")

			}
		}
	}(msgs)
	return srvc
}
