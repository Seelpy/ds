package inmemory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"protokey/pkg/app"
)

type Config struct {
	CommandLogFile   string
	SnapshotFile     string
	SnapshotInterval time.Duration
	FlushInterval    time.Duration
}

type Store struct {
	config          Config
	commandChan     chan app.Command
	responseChan    chan app.Response
	pendingCommands []app.Command
	data            map[string]int
	dataFileHandle  *os.File
	mu              sync.RWMutex
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

func NewStore(config Config, commandChan chan app.Command, responseChan chan app.Response) *Store {
	store := &Store{
		config:       config,
		commandChan:  commandChan,
		responseChan: responseChan,
		data:         make(map[string]int),
		stopChan:     make(chan struct{}),
	}

	if err := store.initialize(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: initialization error: %v\n", err)
	}

	store.wg.Add(1)
	go store.runLoop()

	return store
}

func (s *Store) initialize() error {
	if err := s.loadSnapshot(); err != nil {
		return fmt.Errorf("failed to load snapshot: %w", err)
	}

	if err := s.replayCommandLog(); err != nil {
		return fmt.Errorf("failed to replay command log: %w", err)
	}

	return s.openCommandLog()
}

func (s *Store) loadSnapshot() error {
	filePath := filepath.Clean(s.config.SnapshotFile)
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(data, &s.data)
}

func (s *Store) replayCommandLog() error {
	filePath := filepath.Clean(s.config.CommandLogFile)
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем временный канал для воспроизведения команд
	replayChan := make(chan app.Command, 100)
	defer close(replayChan)

	// Горутина для чтения команд из файла
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var cmd app.Command
			if err := json.Unmarshal(scanner.Bytes(), &cmd); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "warning: invalid command in log: %v\n", err)
				continue
			}
			replayChan <- cmd
		}
	}()

	for cmd := range replayChan {
		s.mu.Lock()
		switch cmd := cmd.(type) {
		case *app.SetCommand:
			s.data[cmd.Key] = cmd.Value
		}
		s.mu.Unlock()
	}

	return nil
}

func (s *Store) openCommandLog() error {
	f, err := os.OpenFile(s.config.CommandLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open command log: %w", err)
	}
	s.dataFileHandle = f
	return nil
}

func (s *Store) runLoop() {
	defer s.wg.Done()

	snapshotTicker := time.NewTicker(s.config.SnapshotInterval)
	defer snapshotTicker.Stop()

	flushTicker := time.NewTicker(s.config.FlushInterval)
	defer flushTicker.Stop()

	for {
		select {
		case cmd, ok := <-s.commandChan:
			if !ok {
				return
			}
			s.processCommand(cmd)

		case <-flushTicker.C:
			s.flushPendingCommands()

		case <-snapshotTicker.C:
			if err := s.createSnapshot(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error creating snapshot: %v\n", err)
			}

		case <-s.stopChan:
			s.flushPendingCommands()
			if err := s.createSnapshot(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error creating final snapshot: %v\n", err)
			}
			return
		}
	}
}

func (s *Store) processCommand(cmd app.Command) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch cmd := cmd.(type) {
	case *app.SetCommand:
		s.data[cmd.Key] = cmd.Value
		s.pendingCommands = append(s.pendingCommands, cmd)
		s.responseChan <- &app.SetResponse{Err: nil}

	case *app.GetCommand:
		value, exists := s.data[cmd.Key]
		if !exists {
			value = 0
		}
		s.responseChan <- &app.GetResponse{Value: value, Err: nil}

	case *app.ListKeysCommand:
		var keys []string
		for k := range s.data {
			if len(k) >= len(cmd.Prefix) && k[:len(cmd.Prefix)] == cmd.Prefix {
				keys = append(keys, k)
			}
		}
		s.responseChan <- &app.ListKeysResponse{Keys: keys, Err: nil}

	default:
		s.responseChan <- &app.SetResponse{Err: fmt.Errorf("unknown operation %d", cmd.GetType())}
	}
}

func (s *Store) flushPendingCommands() {
	if s.dataFileHandle == nil || len(s.pendingCommands) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	writer := bufio.NewWriter(s.dataFileHandle)
	for _, cmd := range s.pendingCommands {
		data, err := json.Marshal(cmd)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshaling command: %v\n", err)
			continue
		}

		if _, err := writer.Write(data); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error writing command: %v\n", err)
			return
		}

		if _, err := writer.WriteString("\n"); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error writing newline: %v\n", err)
			return
		}
	}

	if err := writer.Flush(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error flushing commands: %v\n", err)
	}

	s.pendingCommands = nil
}

func (s *Store) createSnapshot() error {
	s.mu.RLock()
	dataCopy := make(map[string]int, len(s.data))
	for k, v := range s.data {
		dataCopy[k] = v
	}
	s.mu.RUnlock()

	filePath := filepath.Clean(s.config.SnapshotFile)
	tmpFilePath := filePath + ".tmp"

	data, err := json.Marshal(dataCopy)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(tmpFilePath, data, 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return fmt.Errorf("rename error: %w", err)
	}

	return nil
}

func (s *Store) Close() {
	close(s.stopChan)
	s.wg.Wait()

	if s.dataFileHandle != nil {
		_ = s.dataFileHandle.Close()
	}
}
