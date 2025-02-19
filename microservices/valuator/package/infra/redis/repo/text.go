package repo

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"github.com/redis/go-redis/v9"
	"valuator/package/app/model"
	"valuator/package/infra/keyvalue"
)

const (
	keyPrefix = "text:"
	allQuery  = "text:*"
)

func NewTextRepository(client *redis.Client) model.TextRepository {
	return &textRepository{
		storage: keyvalue.NewStorage[textSerializable](client),
	}
}

type textSerializable struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type textRepository struct {
	storage keyvalue.Storage[textSerializable]
}

func (r *textRepository) Store(text model.Text) error {
	return r.storage.Set(context.Background(), keyPrefix+uuid.UUID(text.ID()).String(), textSerializable{
		ID:    uuid.UUID(text.ID()).String(),
		Value: text.Value(),
	}, 0)
}

func (r *textRepository) Create(value string) model.Text {
	return model.NewText(value)
}

func (r *textRepository) Remove(text model.Text) error {
	return r.storage.Delete(context.Background(), keyPrefix+uuid.UUID(text.ID()).String())
}

func (r *textRepository) Find(id model.TextID) (maybe.Maybe[model.Text], error) {
	v, err := r.storage.Get(context.Background(), keyPrefix+uuid.UUID(id).String())
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
	vs, err := r.storage.ListAll(context.Background(), allQuery)
	if err != nil {
		return nil, err
	}
	texts := make([]model.Text, 0, len(vs))
	for _, v := range vs {
		textModel, err1 := r.convertToModel(v)
		if err1 != nil {
			return nil, err1
		}
		texts = append(texts, textModel)
	}
	return texts, nil
}

func (r *textRepository) convertToModel(text textSerializable) (model.Text, error) {
	id, err := uuid.FromString(text.ID)
	if err != nil {
		return nil, err
	}
	return model.LoadText(
		model.TextID(id),
		text.Value,
	), nil
}
