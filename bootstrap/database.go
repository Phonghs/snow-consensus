package bootstrap

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func NewRedisClient(env *Env) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisHost := env.RedisHost
	redisPort := env.RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(redisHost, ":", redisPort),
		Password: env.RedisPassword,
		DB:       0, // use default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Println("Redis can't be connected: ", err)
	}
	return redisClient
}

func CloseRedisClient(redisClient *redis.Client) {
	err := redisClient.Close()
	if err != nil {
		log.Println("Redis can't be closed: ", err)
	}
}
