package unique

import (
	"context"
	"crypto/sha256"
	"errors"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/app/unique"
	"rankcalculator/package/infra/redis/keyvalue"
)

const (
	keyPrefix = "unique:"
)

func NewUniqueStorage(client *redis.Client) unique.TextUniquenessStorage {
	return &uniqueStorage{
		storage: keyvalue.NewStorage[string](client),
	}
}

type uniqueStorage struct {
	storage keyvalue.Storage[string]
}

func (r *uniqueStorage) IsUnique(text string) (bool, error) {
	_, err := r.storage.Get(context.Background(), keyPrefix+hash(text))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

func (r *uniqueStorage) Store(text string) error {
	keyHash := hash(text)
	return r.storage.Set(context.Background(), keyPrefix+keyHash, "", 0)
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}
