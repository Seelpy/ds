package handler

import (
	"github.com/gofrs/uuid"
	"rankcalculator/package/app/service"
)

func NewRemoveTextStatisticsCommandHandler(service service.StatisticsService) TextStatisticsCommandHandler {
	return &removeTextStatisticsCommandHandler{
		service: service,
	}
}

type removeTextStatisticsCommandHandler struct {
	service service.StatisticsService
}

func (h *removeTextStatisticsCommandHandler) Handle(userID uuid.UUID, textID uuid.UUID, textValue string) error {
	return h.service.RemoveStatistics(userID, textID, textValue)
}
