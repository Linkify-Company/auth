package main

import (
	"auth/internal/config"
	"auth/internal/handler"
	v1 "auth/internal/handler/v1"
	"auth/internal/repository"
	"auth/internal/repository/postgres"
	"auth/internal/repository/redis"
	"auth/internal/service"
	"context"
	"fmt"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.GetLogger(cfg.Application.Env)

	pool, err := postgres.New(context.Background(), log, cfg.Application.Debug)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	redisClient, err := redis.New()
	if err != nil {
		log.Error(err)
		panic(err)
	}

	authService := service.NewService(
		log,
		pool,
		redisClient,
		repository.NewRepository(),
		&cfg.EmailService,
	)

	router := handler.Run(
		log,
		v1.NewHandler(&cfg.Handler, &cfg.Token, log, authService),
	)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port), router); err != nil {
			log.Error(errify.NewInternalServerError(err.Error(), "main/ListenAndServe"))
			panic(err)
		}
	}()

	log.Infof("server listening on port %d", cfg.Server.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Infof("application stopped")
}
