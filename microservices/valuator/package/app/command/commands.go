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
	GetTextID() uuid.UUID
}

func NewCalculateCommand(TextID uuid.UUID) Command {
	return &CalculateCommand{TextID: TextID}
}

type CalculateCommand struct {
	TextID uuid.UUID `json:"textID"`
}

func (c *CalculateCommand) Type() CommandType {
	return CalculateCommandType
}

func (c *CalculateCommand) GetTextID() uuid.UUID {
	return c.TextID
}

func NewRemoveCommand(TextID uuid.UUID) Command {
	return &RemoveCommand{TextID: TextID}
}

type RemoveCommand struct {
	TextID uuid.UUID `json:"textID"`
}

func (c *RemoveCommand) Type() CommandType {
	return RemoveCommandType
}

func (c *RemoveCommand) GetTextID() uuid.UUID {
	return c.TextID
}
