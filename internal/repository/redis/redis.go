package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository interface {
	Banner
}

type Redis struct {
	client *redis.Client
}

const (
	CacheLive = time.Minute * 5
)

func New(client *redis.Client) Repository {
	return &Redis{
		client,
	}
}
