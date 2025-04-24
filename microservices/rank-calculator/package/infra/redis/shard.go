package redis

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/infra/redis/keyvalue"
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
