package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLock struct {
	redis *redis.Client
	retry int
}

type RedisLock interface {
	Acquire(ctx context.Context, key string, expiry, delay time.Duration) bool
	Release(ctx context.Context, key string) error
}

var rdb *redis.Client

func NewRedisLock() RedisLock {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_HOST"),
			DB:   0,
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.Fatalln("[Redis][Init] error while ping client")
		}
	}

	return &redisLock{
		redis: rdb,
		retry: 3,
	}
}

func (c *redisLock) Acquire(ctx context.Context, key string, expiry, delay time.Duration) bool {
	for i := 0; i < c.retry; i++ {
		ok, err := c.redis.SetNX(ctx, key, "locked", expiry).Result()
		if err != nil {
			fmt.Println("[Redis][Acquire] error to acquire lock", err.Error())
			return false
		}
		if ok {
			return true
		}
		time.Sleep(delay)
	}
	fmt.Println("[Redis][Acquire] Failed to acquire lock")
	return false
}

func (c *redisLock) Release(ctx context.Context, key string) error {
	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("[Redis][Release] error while release", err.Error())
	}
	return err
}
