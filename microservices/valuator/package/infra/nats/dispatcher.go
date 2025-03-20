package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"valuator/package/app/command"
)

type NatsDispatcher struct {
	conn *nats.Conn
}

func NewNatsDispatcher(conn *nats.Conn) *NatsDispatcher {
	return &NatsDispatcher{conn: conn}
}

func (d *NatsDispatcher) Publish(command command.Command) error {
	subject := getSubjectForCommand(command)
	data, err := serializeCommand(command)
	if err != nil {
		return err
	}
	return d.conn.Publish(subject, data)
}

func (d *NatsDispatcher) Close() {
	d.conn.Close()
}

func getSubjectForCommand(cmd command.Command) string {
	switch cmd.Type() {
	case command.CalculateCommandType:
		return "command.calculate"
	case command.RemoveCommandType:
		return "command.remove"
	default:
		return "command.unknown"
	}
}

func serializeCommand(command command.Command) ([]byte, error) {
	return json.Marshal(command)
}
