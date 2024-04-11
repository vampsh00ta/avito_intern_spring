package redis

import (
	"avito_intern/internal/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Banner interface {
	GetUserBanner(ctx context.Context, tagID, featureID int32) (models.Banner, error)
	SetUserBanner(ctx context.Context, tagID, featureID int32, banner models.Banner) error
}

func (r Redis) GetUserBanner(ctx context.Context, tagID, featureID int32) (models.Banner, error) {
	key := strconv.Itoa(int(tagID)) + "_" + strconv.Itoa(int(featureID))
	res := r.client.Get(ctx, key)
	switch {
	case errors.Is(res.Err(), redis.Nil):
		return models.Banner{}, nil
	case res.Err() != nil:
		return models.Banner{}, res.Err()
	}
	if res.Err() != nil {
		return models.Banner{}, res.Err()
	}
	var banner models.Banner
	if err := json.Unmarshal([]byte(res.String()), &banner); err != nil {
		return models.Banner{}, err
	}
	return banner, nil
}

func (r Redis) SetUserBanner(ctx context.Context, tagID, featureID int32, banner models.Banner) error {
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
