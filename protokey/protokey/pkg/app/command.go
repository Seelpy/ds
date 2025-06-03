package app

type OperationType int

const (
	SetOperation OperationType = iota
	GetOperation
	ListKeysOperation
)

type Command interface {
	GetType() OperationType
}

type SetCommand struct {
	Key   string
	Value int
}

func (c *SetCommand) GetType() OperationType {
	return SetOperation
}

type GetCommand struct {
	Key string
}

func (c *GetCommand) GetType() OperationType {
	return GetOperation
}

type ListKeysCommand struct {
	Prefix string
}

func (c *ListKeysCommand) GetType() OperationType {
	return ListKeysOperation
}

type Response interface {
	GetType() OperationType
}

type SetResponse struct {
	Err error
}

func (r *SetResponse) GetType() OperationType {
	return SetOperation
}

type GetResponse struct {
	Value int
	Err   error
}

func (r *GetResponse) GetType() OperationType {
	return GetOperation
}

type ListKeysResponse struct {
	Keys []string
	Err  error
}

func (r *ListKeysResponse) GetType() OperationType {
	return ListKeysOperation
}
