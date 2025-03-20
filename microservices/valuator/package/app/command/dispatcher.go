package command

type Dispatcher interface {
	Publish(command Command) error
}
