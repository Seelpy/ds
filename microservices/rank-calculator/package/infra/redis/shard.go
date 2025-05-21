package redis

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/app/model"
	"rankcalculator/package/infra/redis/keyvalue"
)

type Region string

const (
	ruRegion   Region = "RU"
	euRegion   Region = "EU"
	asiaRegion Region = "ASIA"
)

var countryToRegion = map[model.Country]Region{
	model.RussiaCountry:  ruRegion,
	model.FranceCountry:  euRegion,
	model.GermanyCountry: euRegion,
	model.UAECountry:     asiaRegion,
	model.IndiaCountry:   asiaRegion,
}

const (
	keyPrefix = "user:"
)

type userSerializable struct {
	ID      string `json:"id"`
	Country string `json:"country"`
}

type Provider interface {
	GetRedisShard(userID uuid.UUID) (*redis.Client, error)
	ListAllShards() []*redis.Client
}

func NewShardRepository(client, ruClient, euClient, asiaClient *redis.Client) *TextShardRepository {
	return &TextShardRepository{
		storage: keyvalue.NewStorage[userSerializable](client),
		regionMap: map[Region]*redis.Client{
			ruRegion:   ruClient,
			euRegion:   euClient,
			asiaRegion: asiaClient,
		},
	}
}

type TextShardRepository struct {
	storage   keyvalue.Storage[userSerializable]
	regionMap map[Region]*redis.Client
}

func (r *TextShardRepository) GetRedisShard(userID uuid.UUID) (*redis.Client, error) {
	v, err := r.storage.Get(context.Background(), keyPrefix+userID.String())
	if err != nil {
		return nil, err
	}
	return r.regionMap[countryToRegion[model.Country(v.Country)]], nil
}

func (r *TextShardRepository) ListAllShards() []*redis.Client {
	result := make([]*redis.Client, 0)
	for _, client := range r.regionMap {
		result = append(result, client)
	}
	return result
}
