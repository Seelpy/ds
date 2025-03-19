package nats

type Handler struct {
}

func (h *Handler) Handle(msgs [][]byte) error {
	return nil
}
