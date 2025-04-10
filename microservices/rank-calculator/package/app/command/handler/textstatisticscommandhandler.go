package handler

import "github.com/gofrs/uuid"

type TextStatisticsCommandHandler interface {
	Handle(textID uuid.UUID, textValue string) error
}
