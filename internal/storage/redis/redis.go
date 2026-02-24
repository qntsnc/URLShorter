package redis

import (
	"context"

	"log"
	"os"

	rs "github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client *rs.Client
}

func NewRedisStorage() (*RedisStorage, error) {
	conn := rs.NewClient(&rs.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err := conn.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Unable to connect to Redis")
		return &RedisStorage{}, err
	}
	return &RedisStorage{Client: conn}, nil
}
