package infra

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"protokey/pkg/app"
	"sync"
)

type jsonCommand struct {
	Type  app.OperationType `json:"type"`
	Key   string            `json:"key"`
	Value int               `json:"value"`
}

type FileSnapshotService struct {
	filePath string
	mu       sync.Mutex
	file     *os.File
	writer   *bufio.Writer
}

func NewFileSnapshotService(filePath string) (*FileSnapshotService, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open snapshot file: %w", err)
	}

	return &FileSnapshotService{
		filePath: filePath,
		file:     file,
		writer:   bufio.NewWriter(file),
	}, nil
}

func (fss *FileSnapshotService) Close() error {
	fss.mu.Lock()
	defer fss.mu.Unlock()

	if fss.writer != nil {
		if err := fss.writer.Flush(); err != nil {
			fss.file.Close()
			return err
		}
	}
	if fss.file != nil {
		err := fss.file.Close()
		fss.file = nil
		fss.writer = nil
		return err
	}
	return nil
}

func (fss *FileSnapshotService) Append(commands []app.Command) error {
	if len(commands) == 0 {
		return nil
	}

	fss.mu.Lock()
	defer fss.mu.Unlock()

	for _, cmd := range commands {
		if setCmd, ok := cmd.(*app.SetCommand); ok {
			jsonCmd := jsonCommand{
				Type:  app.SetOperation,
				Key:   setCmd.Key,
				Value: setCmd.Value,
			}

			data, err := json.Marshal(jsonCmd)
			if err != nil {
				return fmt.Errorf("failed to marshal command: %w", err)
			}

			if _, err := fss.writer.Write(data); err != nil {
				return fmt.Errorf("failed to write command: %w", err)
			}

			if _, err := fss.writer.WriteString("\n"); err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}
	}

	return fss.writer.Flush()
}

func (fss *FileSnapshotService) GetSnapshot() ([]app.Command, error) {
	fss.mu.Lock()
	defer fss.mu.Unlock()

	if fss.file != nil {
		if err := fss.writer.Flush(); err != nil {
			fss.file.Close()
			return nil, err
		}
		fss.file.Close()
		fss.file = nil
		fss.writer = nil
	}

	file, err := os.Open(fss.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []app.Command{}, nil
		}
		return nil, fmt.Errorf("failed to open snapshot file: %w", err)
	}
	defer file.Close()

	var commands []app.Command
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var cmd jsonCommand
		if err := json.Unmarshal(scanner.Bytes(), &cmd); err != nil {
			return nil, fmt.Errorf("failed to unmarshal command: %w", err)
		}

		commands = append(commands, &app.SetCommand{
			Key:   cmd.Key,
			Value: cmd.Value,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fss.file, err = os.OpenFile(fss.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen file for writing: %w", err)
	}
	fss.writer = bufio.NewWriter(fss.file)

	return commands, nil
}
