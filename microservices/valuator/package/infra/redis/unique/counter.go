package unique

import (
	"context"
	"crypto/sha256"
	"errors"
	"github.com/redis/go-redis/v9"
	"valuator/package/app/unique"
	"valuator/package/infra/keyvalue"
)

const (
	keyPrefix = "unique:"
)

func NewUniqueCounter(client *redis.Client) unique.TextCounter {
	return &uniqueCounter{
		storage: keyvalue.NewStorage[countSerializable](client),
	}
}

type countSerializable struct {
	Count int `json:"count"`
}

type uniqueCounter struct {
	storage keyvalue.Storage[countSerializable]
}

func (r *uniqueCounter) GetCount(key string) (int, error) {
	result, err := r.storage.Get(context.Background(), keyPrefix+hash(key))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return result.Count, nil
}

func (r *uniqueCounter) Dec(key string) error {
	keyHash := hash(key)
	result, err := r.storage.Get(context.Background(), keyPrefix+keyHash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	if result.Count > 0 {
		result.Count = result.Count - 1
	}
	return r.storage.Set(context.Background(), keyPrefix+keyHash, result, 0)
}

func (r *uniqueCounter) Inc(key string) error {
	keyHash := hash(key)
	result, err := r.storage.Get(context.Background(), keyPrefix+keyHash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return r.storage.Set(context.Background(), keyPrefix+keyHash, countSerializable{Count: 1}, 0)
		}
		return err
	}
	result.Count++
	return r.storage.Set(context.Background(), keyPrefix+keyHash, result, 0)
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}
