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
	GetTextValue() string
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

func (c *CalculateCommand) GetTextValue() string {
	return ""
}

type RemoveCommand struct {
	TextID    uuid.UUID `json:"textID"`
	TextValue string    `json:"textValue"`
}

func (c *RemoveCommand) Type() CommandType {
	return RemoveCommandType
}

func (c *RemoveCommand) GetTextID() uuid.UUID {
	return c.TextID
}

func (c *RemoveCommand) GetTextValue() string {
	return c.TextValue
}
