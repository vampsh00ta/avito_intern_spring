package vk

import (
	"avito_intern/config"
	psqlrep "avito_intern/internal/repository/psql"
	redisrep "avito_intern/internal/repository/redis"
	"avito_intern/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	transport "avito_intern/internal/transport/http"
	"avito_intern/pkg/client"
	"context"
	"fmt"

	//"go.uber.org/zap"

	"log"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {

	ctx := context.Background()

	pg, err := client.New(ctx, 5, cfg.PG)
	if err != nil {
		//log.Fatal(fmt.Errorf("vk - Run - postgres.New: %w", errs))
	}
	psqlrepo := psqlrep.New(pg)

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})
	redisrepo := redisrep.New(clientRedis)

	if err != nil {
		log.Fatal(fmt.Errorf("avito - Run - postgres.New: %w", err))
	}

	srvc := service.New(psqlrepo, redisrepo, cfg)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	t := transport.New(srvc, sugar)

	log.Print("Listening...")

	if err := http.ListenAndServe(":"+cfg.HTTP.Port, t); err != nil {
		panic(err)
	}
	select {
	case <-interrupt:
		panic("exit")
	}

}
