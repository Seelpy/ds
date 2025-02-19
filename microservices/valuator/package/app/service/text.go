package service

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/model"
)

type TextService interface {
	Add(value string) (uuid.UUID, error)
	Remove(id uuid.UUID) error
}

func NewTextService(repo model.TextRepository) TextService {
	return &textService{repo: repo}
}

type textService struct {
	repo model.TextRepository
}

func (s *textService) Add(value string) (uuid.UUID, error) {
	text := s.repo.Create(value)
	err := s.repo.Store(text)
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
		return s.repo.Remove(text.Value())
	}
	return nil
}
