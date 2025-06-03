package infra

import (
	"fmt"
	"protokey/pkg/app"
	"sync"
)

type Store struct {
	data                map[string]int
	processedCommands   []app.Command
	mu                  sync.RWMutex
	processedCommandsMu sync.Mutex
}

func NewStore() *Store {
	return &Store{
		data:              make(map[string]int),
		processedCommands: make([]app.Command, 0),
	}
}

func (s *Store) ProcessCommand(cmd app.Command) app.Response {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.processedCommandsMu.Lock()
	s.processedCommands = append(s.processedCommands, cmd)
	s.processedCommandsMu.Unlock()

	switch cmd := cmd.(type) {
	case *app.SetCommand:
		s.data[cmd.Key] = cmd.Value
		return &app.SetResponse{Err: nil}

	case *app.GetCommand:
		value, exists := s.data[cmd.Key]
		if !exists {
			value = 0
		}
		return &app.GetResponse{Value: value, Err: nil}

	case *app.ListKeysCommand:
		var keys []string
		for k := range s.data {
			if len(k) >= len(cmd.Prefix) && k[:len(cmd.Prefix)] == cmd.Prefix {
				keys = append(keys, k)
			}
		}
		return &app.ListKeysResponse{Keys: keys, Err: nil}

	default:
		return &app.SetResponse{Err: fmt.Errorf("unknown operation %d", cmd.GetType())}
	}
}

func (s *Store) Load(commands []app.Command) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]int)

	for _, cmd := range commands {
		if v, ok := cmd.(*app.SetCommand); ok {
			s.data[v.Key] = v.Value
		}
	}
	return nil
}

func (s *Store) PopProcessedCommands() []app.Command {
	s.processedCommandsMu.Lock()
	defer s.processedCommandsMu.Unlock()

	commands := s.processedCommands
	s.processedCommands = make([]app.Command, 0)
	return commands
}
