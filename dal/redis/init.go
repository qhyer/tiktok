package redis

import (
	"context"

	"tiktok/pkg/constants"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() {
	client := redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddress,
		Password: constants.RedisPassword,
		DB:       0, // use default DB
		PoolSize: 100,
	})

	// 检查是否成功连接到redis服务器
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("init redis failed")
	}

	RDB = client
}
