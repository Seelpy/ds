package notification

import "github.com/gofrs/uuid"

const (
	prefix = "rank"
)

func GenerateChannel(textID uuid.UUID) string {
	return prefix + "#" + textID.String()
}
