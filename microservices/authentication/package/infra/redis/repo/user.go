package repo

import (
	"authentication/package/app/model"
	"authentication/package/infra/keyvalue"
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	userKeyPrefix = "user:"
	loginIndex    = "login_index:"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateLogin = errors.New("login already exists")
)

func NewUserRepository(client *redis.Client) model.UserRepository {
	return &UserRepository{
		keyvalue: keyvalue.NewStorage[userSerializable](client),
		client:   client,
	}
}

type userSerializable struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	Country      string `json:"country"`
}

type UserRepository struct {
	keyvalue keyvalue.Storage[userSerializable]
	client   *redis.Client
}

func (r *UserRepository) Register(input model.UserRegisterInput) (model.User, error) {
	ctx := context.Background()

	// Check if login already exists
	exists, err := r.client.Exists(ctx, loginIndex+input.Login).Result()
	if err != nil {
		return model.User{}, err
	}
	if exists > 0 {
		return model.User{}, ErrDuplicateLogin
	}

	id := uuid.NewV1()

	// Create user object
	user := userSerializable{
		ID:           id.String(),
		Login:        input.Login,
		PasswordHash: input.PasswordHash,
		Country:      string(input.Country),
	}

	// Store user data
	userKey := userKeyPrefix + id.String()
	err = r.keyvalue.Set(ctx, userKey, user, 0)
	if err != nil {
		return model.User{}, err
	}

	// Create login index
	err = r.client.Set(ctx, loginIndex+input.Login, userKey, 0).Err()
	if err != nil {
		// Rollback user creation if index fails
		r.client.Del(ctx, userKey)
		return model.User{}, err
	}

	return model.User{
		UserID:       id,
		Login:        input.Login,
		PasswordHash: input.PasswordHash,
		Country:      input.Country,
	}, nil
}

func (r *UserRepository) GetByID(userID uuid.UUID) (model.User, error) {
	ctx := context.Background()
	userKey := userKeyPrefix + userID.String()

	user, err := r.keyvalue.Get(ctx, userKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	return model.User{
		UserID:       userID,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		Country:      model.Country(user.Country),
	}, nil
}

func (r *UserRepository) GetByLogin(login string) (model.User, error) {
	ctx := context.Background()

	userKey, err := r.client.Get(ctx, loginIndex+login).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	user, err := r.keyvalue.Get(ctx, userKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	id, err := uuid.FromString(user.ID)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		UserID:       id,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		Country:      model.Country(user.Country),
	}, nil
}
