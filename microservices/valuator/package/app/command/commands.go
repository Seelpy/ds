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

func NewCalculateCommand(TextID uuid.UUID, userID uuid.UUID) Command {
	return &CalculateCommand{TextID: TextID, UserID: userID}
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

func NewRemoveCommand(textID uuid.UUID, textValue string, userID uuid.UUID) Command {
	return &RemoveCommand{TextID: textID, TextValue: textValue, UserID: userID}
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
