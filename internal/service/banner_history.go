package service

import (
	"avito_intern/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type BannerHistory interface {
	GetBannerWithHistory(ctx context.Context, bannerID, limit int) ([]models.Banner, error)
	BannerHistoryCleaner(serviceError chan<- error, done <-chan bool, limit int)
}

const (
	bannerHistoryCleanerPeriod = time.Minute * 1
)

func (s service) GetBannerWithHistory(ctx context.Context, bannerID, limit int) ([]models.Banner, error) {
	var res []models.Banner
	if limit <= 0 || limit > 3 {
		limit = 3
	}
	res, err := s.db.GetBannerWithHistory(ctx, bannerID, limit)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s service) BannerHistoryCleaner(msgs chan<- error, done <-chan bool, limit int) {
	ticker := time.NewTicker(bannerHistoryCleanerPeriod)
	ctx := context.Background()
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				var msg error
				if err := s.db.CleanBannerHistory(ctx, limit); err != nil && !errors.Is(err, pgx.ErrNoRows) {
					msg = err
				}
				msgs <- msg
			}
		}
	}()
}
