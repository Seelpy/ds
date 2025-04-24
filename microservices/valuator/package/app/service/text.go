package service

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/command"
	"valuator/package/app/model"
)

type TextService interface {
	Add(value string, country string) (uuid.UUID, error)
	Remove(id uuid.UUID) error
}

func NewTextService(repo model.TextRepository, textCountryRepo model.TextCountryRepository, dispatcher command.Dispatcher) TextService {
	return &textService{repo: repo, textCountryRepo: textCountryRepo, dispatcher: dispatcher}
}

type textService struct {
	repo            model.TextRepository
	textCountryRepo model.TextCountryRepository
	dispatcher      command.Dispatcher
}

func (s *textService) Add(value string, country string) (uuid.UUID, error) {
	text := s.repo.Create(value)
	textID := uuid.UUID(text.ID())
	err := s.textCountryRepo.Store(textID, model.Country(country))
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.repo.Store(text)
	if err != nil {
		return uuid.UUID{}, err
	}

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
		return s.dispatcher.Publish(command.NewRemoveCommand(id, text.Value().Value()))
	}
	return nil
}
