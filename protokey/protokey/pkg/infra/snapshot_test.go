package infra

import (
	"bufio"
	"encoding/json"
	"os"
	"protokey/pkg/app"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileSnapshotService(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		assert.NotNil(t, fss)
		assert.Equal(t, tmpFile.Name(), fss.filePath)
		assert.NotNil(t, fss.file)
		assert.NotNil(t, fss.writer)
	})

	t.Run("invalid file path", func(t *testing.T) {
		invalidPath := "/invalid/path/to/file"
		fss, err := NewFileSnapshotService(invalidPath)

		assert.Nil(t, fss)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open snapshot file")
	})
}

func TestFileSnapshotService_Append(t *testing.T) {
	t.Run("append set commands", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		commands := []app.Command{
			&app.SetCommand{Key: "key1", Value: 1},
			&app.SetCommand{Key: "key2", Value: 2},
		}

		err = fss.Append(commands)
		assert.NoError(t, err)

		file, err := os.Open(tmpFile.Name())
		require.NoError(t, err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		require.NoError(t, scanner.Err())
		assert.Len(t, lines, 2)

		var cmd1, cmd2 jsonCommand
		err = json.Unmarshal([]byte(lines[0]), &cmd1)
		assert.NoError(t, err)
		assert.Equal(t, app.SetOperation, cmd1.Type)
		assert.Equal(t, "key1", cmd1.Key)
		assert.Equal(t, 1, cmd1.Value)

		err = json.Unmarshal([]byte(lines[1]), &cmd2)
		assert.NoError(t, err)
		assert.Equal(t, app.SetOperation, cmd2.Type)
		assert.Equal(t, "key2", cmd2.Key)
		assert.Equal(t, 2, cmd2.Value)
	})

	t.Run("append empty commands", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		err = fss.Append([]app.Command{})
		assert.NoError(t, err)
	})

	t.Run("append non-set commands", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		commands := []app.Command{
			&app.GetCommand{Key: "key1"},
			&app.ListKeysCommand{Prefix: "prefix"},
		}

		err = fss.Append(commands)
		assert.NoError(t, err)

		file, err := os.Open(tmpFile.Name())
		require.NoError(t, err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		assert.False(t, scanner.Scan())
		assert.NoError(t, scanner.Err())
	})
}

func TestFileSnapshotService_GetSnapshot(t *testing.T) {
	t.Run("successful read from existing file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Write test data to file
		commands := []jsonCommand{
			{Type: app.SetOperation, Key: "key1", Value: 1},
			{Type: app.SetOperation, Key: "key2", Value: 2},
		}

		file, err := os.OpenFile(tmpFile.Name(), os.O_WRONLY, 0644)
		require.NoError(t, err)
		writer := bufio.NewWriter(file)

		for _, cmd := range commands {
			data, err := json.Marshal(cmd)
			require.NoError(t, err)
			_, err = writer.Write(data)
			require.NoError(t, err)
			_, err = writer.WriteString("\n")
			require.NoError(t, err)
		}
		writer.Flush()
		file.Close()

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		result, err := fss.GetSnapshot()
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		setCmd1, ok := result[0].(*app.SetCommand)
		require.True(t, ok)
		assert.Equal(t, "key1", setCmd1.Key)
		assert.Equal(t, 1, setCmd1.Value)

		setCmd2, ok := result[1].(*app.SetCommand)
		require.True(t, ok)
		assert.Equal(t, "key2", setCmd2.Key)
		assert.Equal(t, 2, setCmd2.Value)
	})

	t.Run("read from empty file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		result, err := fss.GetSnapshot()
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("read from non-existent file", func(t *testing.T) {
		nonExistentFile := "nonexistentfile.json"
		defer os.Remove(nonExistentFile)

		fss, err := NewFileSnapshotService(nonExistentFile)
		require.NoError(t, err)
		defer fss.Close()

		result, err := fss.GetSnapshot()
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("read corrupted file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		err = os.WriteFile(tmpFile.Name(), []byte("invalid json\n"), 0644)
		require.NoError(t, err)

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)
		defer fss.Close()

		_, err = fss.GetSnapshot()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal command")
	})
}

func TestFileSnapshotService_Close(t *testing.T) {
	t.Run("successful close", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)

		err = fss.Close()
		assert.NoError(t, err)
		assert.Nil(t, fss.file)
		assert.Nil(t, fss.writer)
	})

	t.Run("double close", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)

		err = fss.Close()
		assert.NoError(t, err)

		err = fss.Close()
		assert.NoError(t, err)
	})

	t.Run("close with flush error", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err)
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fss, err := NewFileSnapshotService(tmpFile.Name())
		require.NoError(t, err)

		fss.file.Close()

		err = fss.Close()
		assert.Error(t, err)
	})
}
