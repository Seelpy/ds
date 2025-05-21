package query

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/authorization"
	"valuator/package/app/model"
)

type TextQueryService interface {
	List(ctx authorization.Context) ([]TextData, error)
	Get(ctx authorization.Context, id uuid.UUID) (TextData, error)
}

func NewTextQueryService(repo model.TextReadRepository) TextQueryService {
	return &textQueryService{
		repo: repo,
	}
}

type TextData struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Value  string
}

type textQueryService struct {
	repo model.TextReadRepository
}

func (s *textQueryService) List(ctx authorization.Context) ([]TextData, error) {
	texts, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	results := make([]TextData, 0, len(texts))
	for _, text := range texts {
		if text.UserID() != ctx.UserID() {
			continue
		}
		results = append(results, TextData{
			ID:     uuid.UUID(text.ID()),
			UserID: text.UserID(),
			Value:  text.Value(),
		})
	}
	return results, nil
}

func (s *textQueryService) Get(ctx authorization.Context, id uuid.UUID) (TextData, error) {
	text, err := s.repo.Find(ctx, model.TextID(id))
	if err != nil {
		return TextData{}, err
	}

	if text.IsEmpty() {
		return TextData{}, model.ErrTextNotFound
	}

	textValue := text.Value()

	if textValue.UserID() != ctx.UserID() {
		return TextData{}, model.ErrTextPermissionDenied
	}

	return TextData{
		ID:    uuid.UUID(textValue.ID()),
		Value: textValue.Value(),
	}, nil
}
