package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"github.com/redis/go-redis/v9"
	"valuator/package/app/authorization"
	"valuator/package/app/model"
	"valuator/package/infra/keyvalue"
	infraredis "valuator/package/infra/redis"
)

const (
	keyPrefix = "text:"
	allQuery  = "text:*"
)

func NewTextRepository(redisProvider infraredis.Provider) model.TextRepository {
	return &textRepository{
		redisProvider: redisProvider,
	}
}

type textSerializable struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Value  string `json:"value"`
}

type textRepository struct {
	redisProvider infraredis.Provider
}

func (r *textRepository) Store(ctx authorization.Context, text model.Text) error {
	redisClient, err := r.redisProvider.GetRedisShard(ctx)
	if err != nil {
		return err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	return storage.Set(context.Background(), keyPrefix+uuid.UUID(text.ID()).String(), textSerializable{
		ID:     uuid.UUID(text.ID()).String(),
		UserID: text.UserID().String(),
		Value:  text.Value(),
	}, 0)
}

func (r *textRepository) Create(ctx authorization.Context, value string) model.Text {
	return model.NewText(ctx.UserID(), value)
}

func (r *textRepository) Remove(ctx authorization.Context, text model.Text) error {
	redisClient, err := r.redisProvider.GetRedisShard(ctx)
	if err != nil {
		return err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	return storage.Delete(context.Background(), keyPrefix+uuid.UUID(text.ID()).String())
}

func (r *textRepository) Find(ctx authorization.Context, id model.TextID) (maybe.Maybe[model.Text], error) {
	redisClient, err := r.redisProvider.GetRedisShard(ctx)
	if err != nil {
		return maybe.Nothing[model.Text](), err
	}
	storage := keyvalue.NewStorage[textSerializable](redisClient)
	v, err := storage.Get(context.Background(), keyPrefix+uuid.UUID(id).String())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return maybe.Nothing[model.Text](), nil
		}
		return maybe.Nothing[model.Text](), err
	}

	textModel, err := r.convertToModel(v)
	if err != nil {
		return maybe.Nothing[model.Text](), err
	}
	return maybe.Just(textModel), nil
}

func (r *textRepository) ListAll() ([]model.Text, error) {
	redisClients := r.redisProvider.ListAllShards()
	result := make([]model.Text, 0)
	for _, redisClient := range redisClients {
		storage := keyvalue.NewStorage[textSerializable](redisClient)
		vs, err := storage.ListAll(context.Background(), allQuery)
		if err != nil {
			return nil, err
		}
		for _, v := range vs {
			textModel, err1 := r.convertToModel(v)
			if err1 != nil {
				return nil, err1
			}
			result = append(result, textModel)
		}
	}
	return result, nil
}

func (r *textRepository) convertToModel(text textSerializable) (model.Text, error) {
	id, err := uuid.FromString(text.ID)
	if err != nil {
		return nil, err
	}
	fmt.Println("SADC: ", text.UserID)
	userID, err := uuid.FromString(text.UserID)
	if err != nil {
		return nil, err
	}
	return model.LoadText(
		model.TextID(id),
		userID,
		text.Value,
	), nil
}
