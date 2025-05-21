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
	GetUserID() uuid.UUID
	GetTextValue() string
}

type CalculateCommand struct {
	TextID uuid.UUID `json:"textID"`
	UserID uuid.UUID `json:"userID"`
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

func (c *CalculateCommand) GetUserID() uuid.UUID {
	return c.UserID
}

type RemoveCommand struct {
	TextID    uuid.UUID `json:"textID"`
	UserID    uuid.UUID `json:"userID"`
	TextValue string    `json:"textValue"`
}

func (c *RemoveCommand) Type() CommandType {
	return RemoveCommandType
}

func (c *RemoveCommand) GetTextID() uuid.UUID {
	return c.TextID
}

func (c *RemoveCommand) GetUserID() uuid.UUID {
	return c.UserID
}

func (c *RemoveCommand) GetTextValue() string {
	return c.TextValue
}
