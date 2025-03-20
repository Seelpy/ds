package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"rankcalculator/package/app/command"
)

type NATSHandler struct {
	commandHandler *command.Handler
	natsConn       *nats.Conn
}

func NewNATSHandler(natsConn *nats.Conn, commandHandler *command.Handler) *NATSHandler {
	return &NATSHandler{
		commandHandler: commandHandler,
		natsConn:       natsConn,
	}
}

func (h *NATSHandler) Start() error {
	_, err := h.natsConn.Subscribe("command.calculate", func(msg *nats.Msg) {
		log.Println("start handle command.calculate")
		var cmd command.CalculateCommand
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			log.Printf("Failed to unmarshal CalculateCommand: %v", err)
			return
		}
		if err := h.commandHandler.Handle(&cmd); err != nil {
			log.Printf("Failed to handle CalculateCommand: %v", err)
		}
		log.Println("finish handle command.calculate")
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to commands.calculate: %v", err)
	}

	_, err = h.natsConn.Subscribe("command.remove", func(msg *nats.Msg) {
		log.Println("start handle command.remove")
		var cmd command.RemoveCommand
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			log.Printf("Failed to unmarshal RemoveCommand: %v", err)
			return
		}
		if err := h.commandHandler.Handle(&cmd); err != nil {
			log.Printf("Failed to handle RemoveCommand: %v", err)
		}
		log.Println("finish handle command.remove")
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to commands.remove: %v", err)
	}

	return nil
}
