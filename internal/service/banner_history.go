package service

import (
	"avito_intern/internal/models"
	"context"
)

type BannerHistory interface {
	GetBannerWithHistory(ctx context.Context, bannerID, limit int) ([]models.Banner, error)
	bannerHistoryCleaner()
}

func (s service) bannerHistoryCleaner() {

}
