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
	bannerHistoryCleaner(limit int)
}

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

func (s service) bannerHistoryCleaner(limit int) {
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	ctx := context.Background()
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:

				if err := s.db.CleanBannerHistory(ctx, limit); err != nil && !errors.Is(err, pgx.ErrNoRows) {
					<-done
					return
				}

			}
		}
	}()
}
