package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/redis/go-redis/v9"
)

var sizes = []int{30, 1000, 1000000}

func BenchmarkMapSetGet(b *testing.B) {
	for _, size := range sizes {
		b.Run("MapSetGet_"+strconv.Itoa(size), func(b *testing.B) {
			m, keys := PopulateMap(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i%size]
				_ = m[k]
			}
		})
	}
}

func BenchmarkRedisSetGet(b *testing.B) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.FlushDB(ctx)
	for _, size := range sizes {
		b.Run("RedisSetGet_"+strconv.Itoa(size), func(b *testing.B) {
			keys, err := PopulateRedis(ctx, rdb, size)
			if err != nil {
				b.Fatalf("failed to populate redis: %v", err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i%size]
				_, _ = rdb.Get(ctx, k).Result()
			}
		})
	}
}
