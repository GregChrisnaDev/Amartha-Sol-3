package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func Init() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("[Redis][Init] error while ping client")
	}

	return rdb
}

type redisLock struct {
	redis  *redis.Client
	key    string
	expiry time.Duration
	retry  int
	delay  time.Duration
}

type RedisLock interface {
	Acquire(ctx context.Context) bool
	Release(ctx context.Context) error
}

func NewRedisLock(c *redis.Client, key string, expiry, delay time.Duration) RedisLock {
	return &redisLock{
		redis:  c,
		key:    key,
		expiry: expiry,
		retry:  3,
		delay:  delay,
	}
}

func (c *redisLock) Acquire(ctx context.Context) bool {
	for i := 0; i < c.retry; i++ {
		ok, err := c.redis.SetNX(ctx, c.key, "locked", c.expiry).Result()
		if err != nil {
			fmt.Println("[Redis][Acquire] error to acquire lock", err.Error())
			return false
		}
		if ok {
			return true
		}
		time.Sleep(c.delay)
	}
	fmt.Println("[Redis][Acquire] Failed to acquire lock")
	return false
}

func (c *redisLock) Release(ctx context.Context) error {
	err := c.redis.Del(ctx, c.key).Err()
	if err != nil {
		fmt.Println("[Redis][Release] error while release", err.Error())
	}
	return err
}
