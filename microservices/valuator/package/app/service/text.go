package service

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/command"
	"valuator/package/app/model"
	"valuator/package/app/unique"
)

type TextService interface {
	Add(value string) (uuid.UUID, error)
	Remove(id uuid.UUID) error
}

func NewTextService(repo model.TextRepository, counter unique.TextCounter, dispatcher command.Dispatcher) TextService {
	return &textService{repo: repo, counter: counter, dispatcher: dispatcher}
}

type textService struct {
	repo       model.TextRepository
	counter    unique.TextCounter
	dispatcher command.Dispatcher
}

func (s *textService) Add(value string) (uuid.UUID, error) {
	text := s.repo.Create(value)
	err := s.repo.Store(text)
	if err != nil {
		return uuid.UUID{}, err
	}
	textID := uuid.UUID(text.ID())
	err = s.dispatcher.Publish(command.NewCalculateCommand(textID))
	return textID, err
}

func (s *textService) Remove(id uuid.UUID) error {
	text, err := s.repo.Find(model.TextID(id))
	if err != nil {
		return err
	}
	if text.IsPresent() {
		err := s.repo.Remove(text.Value())
		if err != nil {
			return err
		}
	}
	return s.dispatcher.Publish(command.NewRemoveCommand(id))
}
