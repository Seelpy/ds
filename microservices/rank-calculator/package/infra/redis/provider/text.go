package provider

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/app/provider"
	"rankcalculator/package/infra/keyvalue"
)

const (
	keyPrefix = "text:"
)

func NewTextProvider(client *redis.Client) provider.TextProvider {
	return &textProvider{
		storage: keyvalue.NewStorage[textSerializable](client),
	}
}

type textSerializable struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type textProvider struct {
	storage keyvalue.Storage[textSerializable]
}

func (r *textProvider) Get(id uuid.UUID) (provider.TextData, error) {
	v, err := r.storage.Get(context.Background(), keyPrefix+id.String())
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
