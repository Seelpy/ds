package authorization

import (
	"context"
	"github.com/gofrs/uuid"
)

type Context interface {
	UserID() uuid.UUID
	Country() string
}

type appContext struct {
	context.Context
	userID  uuid.UUID
	country string
}

func (c *appContext) UserID() uuid.UUID {
	return c.userID
}

func (c *appContext) Country() string {
	return c.country
}

type ContextKey string

const (
	contextKey ContextKey = "appContext"
)

func NewContext(ctx context.Context, userID uuid.UUID, country string) Context {
	return &appContext{
		Context: ctx,
		userID:  userID,
		country: country,
	}
}

func FromContext(ctx context.Context) (Context, bool) {
	c, ok := ctx.Value(contextKey).(Context)
	return c, ok
}
