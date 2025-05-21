package handler

import "github.com/gofrs/uuid"

type TextStatisticsCommandHandler interface {
	Handle(userID uuid.UUID, textID uuid.UUID, textValue string) error
}
