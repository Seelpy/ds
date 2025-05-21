package repo

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/app/model"
	infraredis "rankcalculator/package/infra/redis"
	"rankcalculator/package/infra/redis/keyvalue"
)

const (
	keyPrefix = "text-statistics:"
)

func NewTextStatisticsRepository(redisProvider infraredis.Provider) model.TextStatisticsRepository {
	return &textStatisticsRepository{
		redisProvider: redisProvider,
	}
}

type textSerializable struct {
	TextID           string `json:"text_id"`
	IsDuplicate      bool   `json:"is_duplicate"`
	AllAlphabetCount int    `json:"all_alphabet_count"`
	AllCount         int    `json:"all_count"`
}

type textStatisticsRepository struct {
	redisProvider infraredis.Provider
}

func (r *textStatisticsRepository) Get(userID uuid.UUID, id uuid.UUID) (model.TextStatistics, error) {
	redisClient, err := r.redisProvider.GetRedisShard(userID)
	if err != nil {
		return model.TextStatistics{}, err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	v, err := storage.Get(context.Background(), keyPrefix+id.String())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.TextStatistics{}, model.ErrStatisticsNotFound
		}
		return model.TextStatistics{}, err
	}

	return model.TextStatistics{
		TextID:           id,
		AllAlphabetCount: v.AllAlphabetCount,
		AllCount:         v.AllCount,
		IsDuplicate:      v.IsDuplicate,
	}, nil
}

func (r *textStatisticsRepository) Store(userID uuid.UUID, textStatistics model.TextStatistics) error {
	redisClient, err := r.redisProvider.GetRedisShard(userID)
	if err != nil {
		return err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	return storage.Set(context.Background(), keyPrefix+textStatistics.TextID.String(), textSerializable{
		TextID:           textStatistics.TextID.String(),
		IsDuplicate:      textStatistics.IsDuplicate,
		AllAlphabetCount: textStatistics.AllAlphabetCount,
		AllCount:         textStatistics.AllCount,
	}, 0)
}

func (r *textStatisticsRepository) Remove(userID uuid.UUID, textID uuid.UUID) error {
	redisClient, err := r.redisProvider.GetRedisShard(userID)
	if err != nil {
		return err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	return storage.Delete(context.Background(), keyPrefix+textID.String())
}
