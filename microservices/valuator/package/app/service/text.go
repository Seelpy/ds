package service

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/model"
	"valuator/package/app/unique"
)

type TextService interface {
	Add(value string) (uuid.UUID, error)
	Remove(id uuid.UUID) error
}

func NewTextService(repo model.TextRepository, counter unique.TextCounter) TextService {
	return &textService{repo: repo, counter: counter}
}

type textService struct {
	repo    model.TextRepository
	counter unique.TextCounter
}

func (s *textService) Add(value string) (uuid.UUID, error) {
	text := s.repo.Create(value)
	err := s.repo.Store(text)
	if err != nil {
		return uuid.UUID{}, err
	}
	err = s.counter.Inc(text.Value())
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.UUID(text.ID()), nil
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
		err = s.counter.Dec(text.Value().Value())
		if err != nil {
			return err
		}
	}
	return nil
}
