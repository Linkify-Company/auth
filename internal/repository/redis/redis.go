package redis

import (
	"auth/internal/config"
	"auth/pkg/errify"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
)

func New() (*redis.Client, errify.IError) {
	index, err := strconv.ParseInt(os.Getenv(config.RedisDatabaseIndex), 10, 32)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "redis.New/ParseInt")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv(config.RedisHost), os.Getenv(config.RedisPort)),
		Password: os.Getenv(config.RedisPassword),
		DB:       int(index),
	})
	status := client.Ping()
	if status.Err() != nil {
		return nil, errify.NewInternalServerError(status.Err().Error(), "redis.New/Ping")
	}
	return client, nil
}
