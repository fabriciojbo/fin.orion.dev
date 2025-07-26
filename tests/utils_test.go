package tests

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestMockOutput testa o MockOutput
func TestMockOutput(t *testing.T) {
	TestSetup(t)

	t.Run("MockOutput captura stdout", func(t *testing.T) {
		output := &MockOutput{}
		testMessage := "test message"

		output.Stdout.WriteString(testMessage)

		assert.Equal(t, testMessage, output.Stdout.String())
	})

	t.Run("MockOutput captura stderr", func(t *testing.T) {
		output := &MockOutput{}
		testError := "test error"

		output.Stderr.WriteString(testError)

		assert.Equal(t, testError, output.Stderr.String())
	})
}

// TestMockFileSystem testa o MockFileSystem
func TestMockFileSystem(t *testing.T) {
	TestSetup(t)

	t.Run("MockFileSystem WriteFile", func(t *testing.T) {
		mockFS := NewMockFileSystem()
		filename := "test.txt"
		data := []byte("test content")

		mockFS.On("WriteFile", filename, data, os.FileMode(0644)).Return(nil)

		err := mockFS.WriteFile(filename, data, 0644)
		assert.NoError(t, err)
		assert.Equal(t, string(data), mockFS.files[filename])
	})

	t.Run("MockFileSystem ReadFile", func(t *testing.T) {
		mockFS := NewMockFileSystem()
		filename := "test.txt"
		content := "test content"

		mockFS.files[filename] = content
		mockFS.On("ReadFile", filename).Return([]byte(content), nil)

		data, err := mockFS.ReadFile(filename)
		assert.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("MockFileSystem Stat", func(t *testing.T) {
		mockFS := NewMockFileSystem()
		filename := "test.txt"

		// Mock file info
		mockFileInfo := &mockFileInfo{
			name:    filename,
			size:    100,
			mode:    0644,
			isDir:   false,
			modTime: time.Now(),
		}

		mockFS.On("Stat", filename).Return(mockFileInfo, nil)

		fileInfo, err := mockFS.Stat(filename)
		assert.NoError(t, err)
		assert.Equal(t, filename, fileInfo.Name())
	})

	t.Run("MockFileSystem MkdirAll", func(t *testing.T) {
		mockFS := NewMockFileSystem()
		path := "test/dir"

		mockFS.On("MkdirAll", path, os.FileMode(0755)).Return(nil)

		err := mockFS.MkdirAll(path, 0755)
		assert.NoError(t, err)
	})
}

// TestMockDocker testa o MockDocker
func TestMockDocker(t *testing.T) {
	TestSetup(t)

	t.Run("MockDocker Run", func(t *testing.T) {
		mockDocker := NewMockDocker()
		args := []string{"ps", "-a"}

		mockDocker.On("Run", args).Return(nil)

		err := mockDocker.Run(args...)
		assert.NoError(t, err)
	})

	t.Run("MockDocker GetContainers", func(t *testing.T) {
		mockDocker := NewMockDocker()

		// Adicionar containers
		mockDocker.SetContainerStatus("container1", "running")
		mockDocker.SetContainerStatus("container2", "stopped")

		containers := mockDocker.GetContainers()
		assert.Equal(t, "running", containers["container1"])
		assert.Equal(t, "stopped", containers["container2"])
	})

	t.Run("MockDocker SetContainerStatus", func(t *testing.T) {
		mockDocker := NewMockDocker()

		mockDocker.SetContainerStatus("test-container", "running")

		containers := mockDocker.GetContainers()
		assert.Equal(t, "running", containers["test-container"])
	})
}

// TestMockServiceBus testa o MockServiceBus
func TestMockServiceBus(t *testing.T) {
	TestSetup(t)

	t.Run("MockServiceBus SendMessage", func(t *testing.T) {
		mockSB := NewMockServiceBus()
		queue := "test.queue"
		message := `{"test": "message"}`

		mockSB.On("SendMessage", queue, message).Return(nil)

		err := mockSB.SendMessage(queue, message)
		assert.NoError(t, err)
		assert.Contains(t, mockSB.queues[queue], message)
	})

	t.Run("MockServiceBus GetMessages", func(t *testing.T) {
		mockSB := NewMockServiceBus()
		queue := "test.queue"
		messages := []string{"msg1", "msg2"}

		mockSB.queues[queue] = messages
		mockSB.On("GetMessages", queue).Return(messages, nil)

		result, err := mockSB.GetMessages(queue)
		assert.NoError(t, err)
		assert.Equal(t, messages, result)
	})

	t.Run("MockServiceBus GetQueues", func(t *testing.T) {
		mockSB := NewMockServiceBus()

		// Adicionar filas
		mockSB.queues["queue1"] = []string{}
		mockSB.queues["queue2"] = []string{}

		mockSB.On("GetQueues").Return([]string{"queue1", "queue2"}, nil)

		queues, err := mockSB.GetQueues()
		assert.NoError(t, err)
		assert.Len(t, queues, 2)
		assert.Contains(t, queues, "queue1")
		assert.Contains(t, queues, "queue2")
	})

	t.Run("MockServiceBus GetTopics", func(t *testing.T) {
		mockSB := NewMockServiceBus()

		// Adicionar tópicos
		mockSB.topics["topic1"] = []string{}
		mockSB.topics["topic2"] = []string{}

		mockSB.On("GetTopics").Return([]string{"topic1", "topic2"}, nil)

		topics, err := mockSB.GetTopics()
		assert.NoError(t, err)
		assert.Len(t, topics, 2)
		assert.Contains(t, topics, "topic1")
		assert.Contains(t, topics, "topic2")
	})
}

// TestTestData testa o TestData
func TestTestData(t *testing.T) {
	TestSetup(t)

	t.Run("GetTestData retorna dados válidos", func(t *testing.T) {
		testData := GetTestData()

		assert.NotEmpty(t, testData.ValidJSON)
		assert.NotEmpty(t, testData.InvalidJSON)
		assert.NotEmpty(t, testData.TestMessage)
		assert.NotEmpty(t, testData.TestQueue)
		assert.NotEmpty(t, testData.TestTopic)

		// Verificar se o JSON válido é realmente válido
		assert.Contains(t, testData.ValidJSON, `"test": "message"`)
		assert.Contains(t, testData.ValidJSON, `"id": 123`)

		// Verificar se o JSON inválido é realmente inválido
		assert.Contains(t, testData.InvalidJSON, `"test": "message"`)
		assert.NotContains(t, testData.InvalidJSON, `}`) // Falta fechamento
	})
}

// TestAssertions testa as funções de assertion
func TestAssertions(t *testing.T) {
	TestSetup(t)

	t.Run("AssertCommandSuccess com sucesso", func(t *testing.T) {
		stdout := "success output"
		stderr := ""
		var err error

		// Não deve paniquear
		AssertCommandSuccess(t, stdout, stderr, err)
	})

	t.Run("AssertCommandFailure com erro", func(t *testing.T) {
		stdout := ""
		stderr := "error output"
		err := assert.AnError

		// Não deve paniquear
		AssertCommandFailure(t, stdout, stderr, err)
	})

	t.Run("AssertFileExists com arquivo existente", func(t *testing.T) {
		// Criar arquivo temporário
		filename := "testdata/temp.txt"
		err := os.WriteFile(filename, []byte("test"), 0644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(filename) }()

		// Não deve paniquear
		AssertFileExists(t, filename)
	})

	t.Run("AssertFileNotExists com arquivo inexistente", func(t *testing.T) {
		filename := "testdata/nonexistent.txt"

		// Não deve paniquear
		AssertFileNotExists(t, filename)
	})
}

// TestExecuteCommand testa a função executeCommand
func TestExecuteCommand(t *testing.T) {
	TestSetup(t)

	t.Run("executeCommand com comando simples", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, cmd)

		assert.NoError(t, err)
		assert.Empty(t, stderr)
		assert.Empty(t, stdout) // Comando não produz saída
	})

	t.Run("executeCommand com comando que produz saída", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("test output")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, cmd)

		assert.NoError(t, err)
		assert.Empty(t, stderr)
		assert.Contains(t, stdout, "test output")
	})

	t.Run("executeCommand com comando que retorna erro", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
			RunE: func(cmd *cobra.Command, args []string) error {
				return assert.AnError
			},
		}

		_, _, err := executeCommand(t, cmd)

		assert.Error(t, err)
		// Remover verificações de stdout e stderr vazios pois podem conter mensagens de erro
	})
}

// mockFileInfo implementa os.FileInfo para testes
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	isDir   bool
	modTime time.Time
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) Sys() interface{}   { return nil }
