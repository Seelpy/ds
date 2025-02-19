package query

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/model"
)

type TextQueryService interface {
	List() ([]TextData, error)
	Get(id uuid.UUID) (TextData, error)
}

func NewTextQueryService(repo model.TextReadRepository) TextQueryService {
	return &textQueryService{
		repo: repo,
	}
}

type TextData struct {
	ID    uuid.UUID
	Value string
}

type textQueryService struct {
	repo model.TextReadRepository
}

func (s *textQueryService) List() ([]TextData, error) {
	texts, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	results := make([]TextData, 0, len(texts))
	for _, text := range texts {
		results = append(results, TextData{
			ID:    uuid.UUID(text.ID()),
			Value: text.Value(),
		})
	}
	return results, nil
}

func (s *textQueryService) Get(id uuid.UUID) (TextData, error) {
	text, err := s.repo.Find(model.TextID(id))
	if err != nil {
		return TextData{}, err
	}

	if text.IsEmpty() {
		return TextData{}, model.ErrTextNotFound
	}

	textValue := text.Value()

	return TextData{
		ID:    uuid.UUID(textValue.ID()),
		Value: textValue.Value(),
	}, nil
}
