package command

import (
	"github.com/gofrs/uuid"
)

type CommandType int

const (
	CalculateCommandType = CommandType(iota)
	RemoveCommandType
)

type Command interface {
	Type() CommandType
	TextID() uuid.UUID
}

type CalculateCommand struct {
	textID uuid.UUID
}

func (c *CalculateCommand) Type() CommandType {
	return CalculateCommandType
}

func (c *CalculateCommand) TextID() uuid.UUID {
	return c.textID
}

type RemoveCommand struct {
	commandType CommandType
	textID      uuid.UUID
}

func (c *RemoveCommand) Type() CommandType {
	return RemoveCommandType
}

func (c *RemoveCommand) TextID() uuid.UUID {
	return c.textID
}
