package vk

import (
	"avito_intern/config"
	psqlrep "avito_intern/internal/repository/psql"
	redisrep "avito_intern/internal/repository/redis"
	"github.com/redis/go-redis/v9"

	"avito_intern/internal/service"
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
		//log.Fatal(fmt.Errorf("vk - Run - postgres.New: %w", err))
	}
	psqlrepo := psqlrep.New(pg)

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})
	redisrepo := redisrep.New(clientRedis)

	if err != nil {
		log.Fatal(fmt.Errorf("vk - Run - postgres.New: %w", err))
	}

	srvc := service.New(psqlrepo, redisrepo)
	banner, err := srvc.GetBannerForUser(ctx, false, 1, 1)
	fmt.Println(banner, err)
	//
	//// Waiting signal
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	////logger, _ := zap.NewProduction()
	////defer logger.Sync() // flushes buffer, if any
	////sugar := logger.Sugar()
	//t := transport.New(srvc, nil)
	//
	//log.Print("Listening...")
	//
	//if err := http.ListenAndServe(":"+cfg.HTTP.Port, t); err != nil {
	//	panic(err)
	//}
	//select {
	//case <-interrupt:
	//	panic("exit")
	//
	//}
	////
	////// Shutdown
	////err = http.Shutdown()
	////if err != nil {
	////	l.Error(fmt.Errorf("vk - Run - httpServer.Shutdown: %w", err))
	////}

}
