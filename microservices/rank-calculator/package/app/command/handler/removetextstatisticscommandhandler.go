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

func (h *removeTextStatisticsCommandHandler) Handle(textID uuid.UUID) error {
	return h.service.RemoveStatistics(textID)
}
