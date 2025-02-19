package keyvalue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage[T any] interface {
	Set(ctx context.Context, key string, value T, expiration time.Duration) error

	Get(ctx context.Context, key string) (T, error)
	ListAll(ctx context.Context, pattern string) ([]T, error)

	Delete(ctx context.Context, key string) error
}

type redisStorage[T any] struct {
	client *redis.Client
}

func NewStorage[T any](client *redis.Client) Storage[T] {
	return &redisStorage[T]{client: client}
}

func (rs *redisStorage[T]) Set(ctx context.Context, key string, value T, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.client.Set(ctx, key, jsonValue, expiration).Err()
}

func (rs *redisStorage[T]) Get(ctx context.Context, key string) (T, error) {
	var result T
	jsonValue, err := rs.client.Get(ctx, key).Bytes()
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(jsonValue, &result)
	return result, err
}

func (rs *redisStorage[T]) ListAll(ctx context.Context, pattern string) ([]T, error) {
	var results []T

	keys, err := rs.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		var obj T

		jsonValue, err := rs.client.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonValue, &obj)
		if err != nil {
			return nil, err
		}

		results = append(results, obj)
	}

	return results, nil
}

func (rs *redisStorage[T]) Delete(ctx context.Context, key string) error {
	return rs.client.Del(ctx, key).Err()
}
