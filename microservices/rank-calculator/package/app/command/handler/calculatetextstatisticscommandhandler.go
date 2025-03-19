package handler

import (
	"github.com/gofrs/uuid"
	"rankcalculator/package/app/service"
)

func NewCalculateTextStatisticsCommandHandler(service service.StatisticsService) TextStatisticsCommandHandler {
	return &calculateTextStatisticsCommandHandler{
		service: service,
	}
}

type calculateTextStatisticsCommandHandler struct {
	service service.StatisticsService
}

func (h *calculateTextStatisticsCommandHandler) Handle(textID uuid.UUID) error {
	return h.service.RankText(textID)
}
