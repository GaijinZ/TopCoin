package service

import (
	"context"
	"crypto/rand"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func GenerateRandomData(size int) ([]string, []string) {
	keys := make([]string, size)
	values := make([]string, size)
	for i := 0; i < size; i++ {
		keys[i] = random4Bytes() + strconv.Itoa(i)
		values[i] = random4Bytes()
	}
	return keys, values
}

func PopulateMap(size int) (map[string]string, []string) {
	m := make(map[string]string, size)
	keys, values := GenerateRandomData(size)
	for i := 0; i < size; i++ {
		m[keys[i]] = values[i]
	}
	return m, keys
}

func PopulateRedis(ctx context.Context, rdb *redis.Client, size int) ([]string, error) {
	keys, values := GenerateRandomData(size)
	for i := 0; i < size; i++ {
		if err := rdb.Set(ctx, keys[i], values[i], 0).Err(); err != nil {
			return nil, err
		}
	}
	return keys, nil
}

func random4Bytes() string {
	b := make([]byte, 4)
	rand.Read(b)
	return string(b)
}
