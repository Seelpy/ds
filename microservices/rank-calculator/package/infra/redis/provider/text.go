package provider

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/app/provider"
	infraredis "rankcalculator/package/infra/redis"
	"rankcalculator/package/infra/redis/keyvalue"
)

const (
	keyPrefix = "text:"
)

func NewTextProvider(redisProvider infraredis.Provider) provider.TextProvider {
	return &textProvider{
		redisProvider: redisProvider,
	}
}

type textSerializable struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type textProvider struct {
	redisProvider infraredis.Provider
}

func (r *textProvider) Get(userID uuid.UUID, id uuid.UUID) (provider.TextData, error) {
	redisClient, err := r.redisProvider.GetRedisShard(userID)
	if err != nil {
		return provider.TextData{}, err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	v, err := storage.Get(context.Background(), keyPrefix+id.String())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return provider.TextData{}, provider.ErrTextNotFound
		}
		return provider.TextData{}, err
	}

	textData, err := r.convertToApp(v)
	if err != nil {
		return provider.TextData{}, err
	}
	return textData, nil
}

func (r *textProvider) convertToApp(text textSerializable) (provider.TextData, error) {
	id, err := uuid.FromString(text.ID)
	if err != nil {
		return provider.TextData{}, err
	}
	return provider.TextData{
		id,
		text.Value,
	}, nil
}
