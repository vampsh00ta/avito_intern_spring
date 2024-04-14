package app

import (
	"avito_intern/config"
	psqlrep "avito_intern/internal/repository/psql"
	redisrep "avito_intern/internal/repository/redis"
	"avito_intern/internal/service"
	"avito_intern/pkg/client"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	transport "avito_intern/internal/transport/http"
	//"go.uber.org/zap".
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	pg, err := client.New(ctx, 5, cfg.PG)
	if err != nil {
		logger.Fatal("app - Run - postgres.New: %w", zap.Error(err))
	}
	psqlrepo := psqlrep.New(pg)

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	redisrepo := redisrep.New(clientRedis)

	srvcLogger := logger.Sugar()
	srvc := service.New(psqlrepo, redisrepo, cfg, srvcLogger)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)

	trptLogger := logger.Sugar()
	t := transport.New(srvc, trptLogger)

	logger.Info("Listening...")

	if err := http.ListenAndServe(":"+cfg.HTTP.Port, t); err != nil { //nolint:gosec
		logger.Fatal(err.Error())
	}
	select {
	case <-interrupt:
		panic("exit")
	}
}
