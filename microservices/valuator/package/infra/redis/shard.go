package redis

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"valuator/package/app/model"
	"valuator/package/infra/keyvalue"
)

type Region string

const (
	ruRegion   Region = "RU"
	euRegion   Region = "EU"
	asiaRegion Region = "ASIA"
)

const (
	keyPrefix = "shard:"
)

var countryToRegion = map[model.Country]Region{
	model.RussiaCountry:  ruRegion,
	model.FranceCountry:  euRegion,
	model.GermanyCountry: euRegion,
	model.UAECountry:     asiaRegion,
	model.IndiaCountry:   asiaRegion,
}

type Provider interface {
	GetRedisShard(textID uuid.UUID) (*redis.Client, error)
	ListAllShards() []*redis.Client
}

func NewShardRepository(client, ruClient, euClient, asiaClient *redis.Client) *TextShardRepository {
	return &TextShardRepository{
		storage: keyvalue.NewStorage[Region](client),
		regionMap: map[Region]*redis.Client{
			ruRegion:   ruClient,
			euRegion:   euClient,
			asiaRegion: asiaClient,
		},
	}
}

type TextShardRepository struct {
	storage   keyvalue.Storage[Region]
	regionMap map[Region]*redis.Client
}

func (r *TextShardRepository) Store(textID uuid.UUID, country model.Country) error {
	region, ok := countryToRegion[country]
	if !ok {
		return model.ErrCountryNotFound
	}
	return r.storage.Set(context.Background(), keyPrefix+textID.String(), region, 0)
}

func (r *TextShardRepository) GetRedisShard(textID uuid.UUID) (*redis.Client, error) {
	v, err := r.storage.Get(context.Background(), keyPrefix+textID.String())
	if err != nil {
		return nil, err
	}
	return r.regionMap[v], nil
}

func (r *TextShardRepository) ListAllShards() []*redis.Client {
	result := make([]*redis.Client, 0)
	for _, client := range r.regionMap {
		result = append(result, client)
	}
	return result
}
