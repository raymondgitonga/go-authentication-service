package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

func NewClient(ctx context.Context, addr string, pass string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})

	err := pingUntilAvailable(ctx, rdb)

	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func pingUntilAvailable(ctx context.Context, redis *redis.Client) error {
	var err error
	var status string
	for i := 0; i < 10; i++ {
		status, err = redis.Ping(ctx).Result()
		if err == nil {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("%w %s", err, status)
}
