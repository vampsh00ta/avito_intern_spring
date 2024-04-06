package service

import (
	"avito_intern/internal/models"
	"context"
)

type Banner interface {
	GetBannerForUser(ctx context.Context, useLastRevision bool, userTag, featureID int32) (models.Banner, error)
}

func (s service) GetBannerForUser(ctx context.Context, useLastRevision bool, userTag, featureID int32) (models.Banner, error) {
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

	return res, err
}
