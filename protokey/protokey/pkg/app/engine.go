package app

import (
	"log"
	"sync"
	"time"
)

type Store interface {
	ProcessCommand(cmd Command) Response
	Load(commands []Command) error
	PopProcessedCommands() []Command
}

type SnapshotService interface {
	Append(commands []Command) error
	GetSnapshot() ([]Command, error)
}

func NewEngine(store Store, snapshot SnapshotService, commandChan chan Command, responseChan chan Response, stopChan chan struct{}, snapshotDelay time.Duration) (*Engine, error) {
	engine := &Engine{
		store:           store,
		snapshotService: snapshot,
		commandChan:     commandChan,
		responseChan:    responseChan,
		stopChan:        stopChan,
		snapshotDelay:   snapshotDelay,
	}
	err := engine.init()
	if err != nil {
		return nil, err
	}
	return engine, nil
}

type Engine struct {
	store           Store
	snapshotService SnapshotService

	commandChan  chan Command
	responseChan chan Response
	stopChan     chan struct{}

	snapshotDelay time.Duration
	mu            sync.Mutex
}

func (e *Engine) Run() {
	go func() {
		snapshotTicker := time.NewTicker(e.snapshotDelay)
		defer snapshotTicker.Stop()

		for {
			select {
			case cmd, ok := <-e.commandChan:
				if !ok {
					return
				}
				resp := e.store.ProcessCommand(cmd)
				e.responseChan <- resp
			case <-snapshotTicker.C:
				e.flushProcessedCommands()
			case <-e.stopChan:
				e.flushProcessedCommands()
				return
			}
		}
	}()
}

func (e *Engine) init() error {
	commands, err := e.snapshotService.GetSnapshot()
	if err != nil {
		return err
	}
	err = e.store.Load(commands)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) flushProcessedCommands() {
	commands := e.store.PopProcessedCommands()
	if len(commands) == 0 {
		return
	}

	err := e.snapshotService.Append(commands)
	if err != nil {
		log.Println("Cant save commands: ", err)
		return
	}
}
