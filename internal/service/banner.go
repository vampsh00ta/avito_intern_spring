package service

import (
	"avito_intern/internal/models"
	"context"
)

type Banner interface {
	GetBannerForUser(ctx context.Context, userTag, featureID int32, useLastRevision bool) (models.Banner, error)
	GetBanners(ctx context.Context, tagID, featureID, limit, offset int32) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner models.Banner) (int, error)
	DeleteBannerByID(ctx context.Context, ID int) error
}

func (s service) DeleteBannerByID(ctx context.Context, ID int) error {
	if err := s.db.DeleteBannerByID(ctx, ID); err != nil {
		return err
	}

	return nil
}
func (s service) GetBannerForUser(ctx context.Context, userTag, featureID int32, useLastRevision bool) (models.Banner, error) {
	var res models.Banner
	res, err := s.cache.GetUserBanner(ctx, userTag, featureID)
	if err != nil {
		return models.Banner{}, err
	}
	if useLastRevision || res.Id == 0 {
		res, err = s.db.GetBannerForUser(ctx, userTag, featureID)
		if err != nil {
			return models.Banner{}, err
		}
		if err := s.cache.SetUserBanner(ctx, userTag, featureID, res); err != nil {
			return models.Banner{}, err
		}

	}
	if res.IsActive == false {
		res.Content = ""
	}
	return res, err
}
func (s service) GetBanners(ctx context.Context, tagID, featureID, limit, offset int32) ([]models.Banner, error) {
	var res []models.Banner
	res, err := s.db.GetBanners(ctx, tagID, featureID, limit, offset)
	if err != nil {
		return nil, err
	}

	return res, err
}
func (s service) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	txCtx, err := s.db.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer s.db.Commit(ctx)
	res, err := s.db.CreateBanner(txCtx, banner)
	if err != nil {
		defer s.db.Rollback(ctx)

		return -1, err
	}

	return res, err
}
