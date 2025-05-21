package command

import (
	"rankcalculator/package/app/command/handler"
	"rankcalculator/package/app/notification"
	"rankcalculator/package/app/service"
	"time"
)

func NewHandler(service service.StatisticsService, publisher notification.Publisher) *Handler {
	return &Handler{
		calculateCommandHandler: handler.NewCalculateTextStatisticsCommandHandler(service, publisher),
		removeCommandHandler:    handler.NewRemoveTextStatisticsCommandHandler(service),
	}
}

type Handler struct {
	calculateCommandHandler handler.TextStatisticsCommandHandler
	removeCommandHandler    handler.TextStatisticsCommandHandler
}

func (h *Handler) Handle(command Command) error {
	time.Sleep(5 * time.Second)
	switch command.Type() {
	case CalculateCommandType:
		return h.calculateCommandHandler.Handle(command.GetUserID(), command.GetTextID(), command.GetTextValue())
	case RemoveCommandType:
		return h.removeCommandHandler.Handle(command.GetUserID(), command.GetTextID(), command.GetTextValue())
	}
	return nil
}
