package repo

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"rankcalculator/package/app/model"
	"rankcalculator/package/infra/redis/keyvalue"
)

const (
	keyPrefix = "text-statistics:"
)

func NewTextStatisticsRepository(client *redis.Client) model.TextStatisticsRepository {
	return &textStatisticsRepository{
		storage: keyvalue.NewStorage[textSerializable](client),
	}
}

type textSerializable struct {
	TextID           string `json:"text_id"`
	IsDuplicate      bool   `json:"is_duplicate"`
	AllAlphabetCount int    `json:"all_alphabet_count"`
	AllCount         int    `json:"all_count"`
}

type textStatisticsRepository struct {
	storage keyvalue.Storage[textSerializable]
}

func (r *textStatisticsRepository) Get(id uuid.UUID) (model.TextStatistics, error) {
	v, err := r.storage.Get(context.Background(), keyPrefix+id.String())
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

func (r *textStatisticsRepository) Store(textStatistics model.TextStatistics) error {
	return r.storage.Set(context.Background(), keyPrefix+textStatistics.TextID.String(), textSerializable{
		TextID:           textStatistics.TextID.String(),
		IsDuplicate:      textStatistics.IsDuplicate,
		AllAlphabetCount: textStatistics.AllAlphabetCount,
		AllCount:         textStatistics.AllCount,
	}, 0)
}

func (r *textStatisticsRepository) Remove(textID uuid.UUID) error {
	return r.storage.Delete(context.Background(), keyPrefix+textID.String())
}
