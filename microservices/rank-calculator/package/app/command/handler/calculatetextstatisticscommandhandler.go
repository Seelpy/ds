package handler

import (
	"github.com/gofrs/uuid"
	"rankcalculator/package/app/notification"
	"rankcalculator/package/app/service"
	"time"
)

func NewCalculateTextStatisticsCommandHandler(service service.StatisticsService, publisher notification.Publisher) TextStatisticsCommandHandler {
	return &calculateTextStatisticsCommandHandler{
		service:   service,
		publisher: publisher,
	}
}

type calculateTextStatisticsCommandHandler struct {
	service   service.StatisticsService
	publisher notification.Publisher
}

func (h *calculateTextStatisticsCommandHandler) Handle(textID uuid.UUID, _ string) error {
	// Спим чтобы почувстовать задержку
	delay := 3 * time.Second
	time.Sleep(delay)

	text, stat, err := h.service.RankText(textID)
	if err != nil {
		return err
	}

	channel := notification.GenerateChannel(textID)
	return h.publisher.Publish(channel, map[string]interface{}{
		"textID":     textID,
		"rank":       stat.Rank(),
		"similarity": stat.IsDuplicate,
		"textValue":  text.Value,
	})
}
