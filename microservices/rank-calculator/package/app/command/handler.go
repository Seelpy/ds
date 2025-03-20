package command

import (
	"rankcalculator/package/app/command/handler"
	"rankcalculator/package/app/service"
)

func NewHandler(service service.StatisticsService) *Handler {
	return &Handler{
		calculateCommandHandler: handler.NewCalculateTextStatisticsCommandHandler(service),
		removeCommandHandler:    handler.NewRemoveTextStatisticsCommandHandler(service),
	}
}

type Handler struct {
	calculateCommandHandler handler.TextStatisticsCommandHandler
	removeCommandHandler    handler.TextStatisticsCommandHandler
}

func (h *Handler) Handle(command Command) error {
	switch command.Type() {
	case CalculateCommandType:
		return h.calculateCommandHandler.Handle(command.GetTextID())
	case RemoveCommandType:
		return h.removeCommandHandler.Handle(command.GetTextID())
	}
	return nil
}
