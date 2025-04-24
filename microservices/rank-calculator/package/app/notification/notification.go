package notification

type Publisher interface {
	Publish(channel string, data interface{}) error
}
