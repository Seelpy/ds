package infra

import (
	"fmt"
	"protokey/pkg/app"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCommand struct {
}

func (m *mockCommand) GetType() app.OperationType {
	return app.OperationType(-1)
}

func TestNewStore(t *testing.T) {
	store := NewStore()
	assert.NotNil(t, store)
	assert.Empty(t, store.data)
	assert.Empty(t, store.processedCommands)
}

func TestProcessCommand_SetCommand(t *testing.T) {
	store := NewStore()
	cmd := &app.SetCommand{Key: "test", Value: 123}

	resp := store.ProcessCommand(cmd)

	assert.IsType(t, &app.SetResponse{}, resp)
	assert.Nil(t, resp.(*app.SetResponse).Err)
	assert.Equal(t, 123, store.data["test"])
	assert.Len(t, store.processedCommands, 1)
	assert.Equal(t, cmd, store.processedCommands[0])
}

func TestProcessCommand_GetCommand_ExistingKey(t *testing.T) {
	store := NewStore()
	store.data["exist"] = 456
	cmd := &app.GetCommand{Key: "exist"}

	resp := store.ProcessCommand(cmd)

	assert.IsType(t, &app.GetResponse{}, resp)
	assert.Nil(t, resp.(*app.GetResponse).Err)
	assert.Equal(t, 456, resp.(*app.GetResponse).Value)
}

func TestProcessCommand_GetCommand_NonExistingKey(t *testing.T) {
	store := NewStore()
	cmd := &app.GetCommand{Key: "nonexist"}

	resp := store.ProcessCommand(cmd)

	assert.IsType(t, &app.GetResponse{}, resp)
	assert.Nil(t, resp.(*app.GetResponse).Err)
	assert.Equal(t, 0, resp.(*app.GetResponse).Value)
}

func TestProcessCommand_ListKeysCommand(t *testing.T) {
	store := NewStore()
	store.data["user:1"] = 1
	store.data["user:2"] = 2
	store.data["config:1"] = 3
	cmd := &app.ListKeysCommand{Prefix: "user:"}

	resp := store.ProcessCommand(cmd)

	assert.IsType(t, &app.ListKeysResponse{}, resp)
	assert.Nil(t, resp.(*app.ListKeysResponse).Err)
	assert.ElementsMatch(t, []string{"user:1", "user:2"}, resp.(*app.ListKeysResponse).Keys)
}

func TestProcessCommand_ListKeysCommand_EmptyPrefix(t *testing.T) {
	store := NewStore()
	store.data["a"] = 1
	store.data["b"] = 2
	cmd := &app.ListKeysCommand{Prefix: ""}

	resp := store.ProcessCommand(cmd)

	assert.IsType(t, &app.ListKeysResponse{}, resp)
	assert.Len(t, resp.(*app.ListKeysResponse).Keys, 2)
}

func TestProcessCommand_UnknownCommand(t *testing.T) {
	store := NewStore()
	cmd := &mockCommand{}

	resp := store.ProcessCommand(cmd)

	assert.IsType(t, &app.SetResponse{}, resp)
	assert.Error(t, resp.(*app.SetResponse).Err)
}

func TestLoad(t *testing.T) {
	store := NewStore()
	commands := []app.Command{
		&app.SetCommand{Key: "a", Value: 1},
		&app.SetCommand{Key: "b", Value: 2},
	}

	err := store.Load(commands)

	assert.Nil(t, err)
	assert.Equal(t, 1, store.data["a"])
	assert.Equal(t, 2, store.data["b"])
	assert.Len(t, store.data, 2)
}

func TestLoad_EmptyCommands(t *testing.T) {
	store := NewStore()
	store.data["old"] = 99 // Should be cleared

	err := store.Load([]app.Command{})

	assert.Nil(t, err)
	assert.Empty(t, store.data)
}

func TestLoad_NonSetCommands(t *testing.T) {
	store := NewStore()
	commands := []app.Command{
		&app.GetCommand{Key: "a"},
		&app.ListKeysCommand{Prefix: "b"},
	}

	err := store.Load(commands)

	assert.Nil(t, err)
	assert.Empty(t, store.data)
}

func TestPopProcessedCommands(t *testing.T) {
	store := NewStore()
	cmd1 := &app.SetCommand{Key: "a", Value: 1}
	cmd2 := &app.GetCommand{Key: "a"}
	store.ProcessCommand(cmd1)
	store.ProcessCommand(cmd2)

	commands := store.PopProcessedCommands()

	assert.Len(t, commands, 2)
	assert.Equal(t, cmd1, commands[0])
	assert.Equal(t, cmd2, commands[1])
	assert.Empty(t, store.processedCommands)
}

func TestPopProcessedCommands_Empty(t *testing.T) {
	store := NewStore()

	commands := store.PopProcessedCommands()

	assert.Empty(t, commands)
}

func TestConcurrentAccess(t *testing.T) {
	store := NewStore()
	numRoutines := 1000
	var wg sync.WaitGroup
	wg.Add(numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			store.ProcessCommand(&app.SetCommand{Key: key, Value: i})
			resp := store.ProcessCommand(&app.GetCommand{Key: key})
			assert.Equal(t, i, resp.(*app.GetResponse).Value)
		}(i)
	}

	wg.Wait()
	assert.Len(t, store.data, numRoutines)
}
