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
	ChangeBanner(ctx context.Context, ID int, banner models.BannerChange) error
	DeleteBannerByTagAndFeature(ctx context.Context, featureID, tagID int32) (int, error)
}

func (s service) ChangeBanner(ctx context.Context, ID int, banner models.BannerChange) error {
	ctx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer s.db.Commit(ctx)
	currBanner, err := s.db.GetBannerByID(ctx, ID)
	if err != nil {
		s.db.Rollback(ctx)

		return err
	}
	if err := s.db.ChangeBanner(ctx, ID, banner); err != nil {
		s.db.Rollback(ctx)

		return err
	}
	if err := s.db.CreateHistoryBanner(ctx, currBanner); err != nil {
		s.db.Rollback(ctx)
		return err
	}

	return nil
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
	if useLastRevision || res.ID == 0 {
		res, err = s.db.GetBannerForUser(ctx, userTag, featureID)
		if err != nil {
			return models.Banner{}, err
		}
		if err := s.cache.SetUserBanner(ctx, userTag, featureID, res); err != nil {
			return models.Banner{}, err
		}

	}
	if !res.IsActive {
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
	ctx, err := s.db.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer s.db.Commit(ctx)
	res, err := s.db.CreateBanner(ctx, banner)
	if err != nil {
		s.db.Rollback(ctx)

		return -1, err
	}

	return res, err
}

func (s service) DeleteBannerByTagAndFeature(ctx context.Context, featureID, tagID int32) (int, error) {
	ctx, err := s.db.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer s.db.Commit(ctx)
	id, err := s.db.DeleteBannerByTagAndFeature(ctx, featureID, tagID)
	if err != nil {

		s.db.Rollback(ctx)
		return -1, err
	}
	return id, nil
}
