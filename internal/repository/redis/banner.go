package redis

import (
	"avito_intern/internal/models"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Banner interface {
	GetUserBanner(ctx context.Context, tagID, featureID int32) (models.Banner, error)
	SetUserBanner(ctx context.Context, tagID, featureID int32, banner models.Banner) error
}

func (r Redis) GetUserBanner(ctx context.Context, tagID, featureID int32) (models.Banner, error) {
	key := strconv.Itoa(int(tagID)) + "_" + strconv.Itoa(int(featureID))
	res, err := r.client.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return models.Banner{}, nil
	case err != nil:
		return models.Banner{}, err
	}

	var banner models.Banner
	if err := json.Unmarshal([]byte(res), &banner); err != nil {
		return models.Banner{}, err
	}
	return banner, nil
}

func (r Redis) SetUserBanner(ctx context.Context, tagID, featureID int32, banner models.Banner) error {
	banner.Feature = featureID
	bannerBytes, err := json.Marshal(banner)
	if err != nil {
		return err
	}
	key := strconv.Itoa(int(tagID)) + "_" + strconv.Itoa(int(featureID))
	if err := r.client.Set(ctx, key, string(bannerBytes), CacheLive); err != nil {
		return nil
	}
	return err
}
