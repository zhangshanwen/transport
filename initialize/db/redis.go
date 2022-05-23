package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}

var R *redis.Client

func InitRedis(r Redis) {
	logrus.Info("--------init_redis_client_start---------")
	host := r.Host
	port := r.Port
	if host == "" {
		host = "localhost"
	}
	if port <= 0 {
		port = 6379
	}
	R = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", host, port),
		Password: r.Password,
		DB:       r.Database,
	})
	if err := R.Ping(context.Background()).Err(); err != nil {
		logrus.Panic(err)
	}
	logrus.Info("--------init_redis_client_end---------")
}
