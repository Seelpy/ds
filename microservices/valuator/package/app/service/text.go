package service

import (
	"github.com/gofrs/uuid"
	"valuator/package/app/authorization"
	"valuator/package/app/command"
	"valuator/package/app/model"
)

type TextService interface {
	Add(ctx authorization.Context, value string) (uuid.UUID, error)
	Remove(ctx authorization.Context, id uuid.UUID) error
}

func NewTextService(repo model.TextRepository, dispatcher command.Dispatcher) TextService {
	return &textService{repo: repo, dispatcher: dispatcher}
}

type textService struct {
	repo       model.TextRepository
	dispatcher command.Dispatcher
}

func (s *textService) Add(ctx authorization.Context, value string) (uuid.UUID, error) {
	text := s.repo.Create(ctx, value)
	textID := uuid.UUID(text.ID())

	err := s.repo.Store(ctx, text)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.dispatcher.Publish(command.NewCalculateCommand(textID, ctx.UserID()))
	return textID, err
}

func (s *textService) Remove(ctx authorization.Context, id uuid.UUID) error {
	text, err := s.repo.Find(ctx, model.TextID(id))
	if err != nil {
		return err
	}

	if text.IsPresent() {
		textValue := text.Value()
		if textValue.UserID() != ctx.UserID() {
			return model.ErrTextPermissionDenied
		}
		err := s.repo.Remove(ctx, text.Value())
		if err != nil {
			return err
		}
		return s.dispatcher.Publish(command.NewRemoveCommand(id, text.Value().Value(), ctx.UserID()))
	}
	return nil
}
