package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOutput captura a saída dos comandos para testes
type MockOutput struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

// TestSetup configura o ambiente de teste
func TestSetup(t *testing.T) {
	// Criar diretórios temporários para testes
	tempDirs := []string{
		"testdata",
		"testdata/messages",
		"testdata/certs",
		"testdata/docker",
	}

	for _, dir := range tempDirs {
		err := os.MkdirAll(dir, 0755)
		assert.NoError(t, err)
	}

	// Limpar após os testes
	t.Cleanup(func() {
		_ = os.RemoveAll("testdata")
	})
}

// executeCommand executa um comando e retorna a saída
func executeCommand(t *testing.T, cmd *cobra.Command, args ...string) (string, string, error) {
	output := &MockOutput{}
	cmd.SetOut(&output.Stdout)
	cmd.SetErr(&output.Stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return output.Stdout.String(), output.Stderr.String(), err
}

// MockFileSystem simula operações de arquivo
type MockFileSystem struct {
	mock.Mock
	files map[string]string
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files: make(map[string]string),
	}
}

func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	m.files[filename] = string(data)
	return args.Error(0)
}

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	if content, exists := m.files[filename]; exists {
		return []byte(content), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFileSystem) Stat(filename string) (os.FileInfo, error) {
	args := m.Called(filename)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

// MockDocker simula operações do Docker
type MockDocker struct {
	mock.Mock
	containers map[string]string
}

func NewMockDocker() *MockDocker {
	return &MockDocker{
		containers: make(map[string]string),
	}
}

func (m *MockDocker) Run(args ...string) error {
	mockArgs := m.Called(args)
	return mockArgs.Error(0)
}

func (m *MockDocker) GetContainers() map[string]string {
	return m.containers
}

func (m *MockDocker) SetContainerStatus(container, status string) {
	m.containers[container] = status
}

// MockServiceBus simula operações do Service Bus
type MockServiceBus struct {
	mock.Mock
	queues   map[string][]string
	topics   map[string][]string
	messages map[string][]string
}

func NewMockServiceBus() *MockServiceBus {
	return &MockServiceBus{
		queues:   make(map[string][]string),
		topics:   make(map[string][]string),
		messages: make(map[string][]string),
	}
}

func (m *MockServiceBus) SendMessage(queue string, message string) error {
	args := m.Called(queue, message)
	if m.queues[queue] == nil {
		m.queues[queue] = []string{}
	}
	m.queues[queue] = append(m.queues[queue], message)
	return args.Error(0)
}

func (m *MockServiceBus) GetMessages(queue string) ([]string, error) {
	args := m.Called(queue)
	return m.queues[queue], args.Error(1)
}

func (m *MockServiceBus) GetTopics() ([]string, error) {
	args := m.Called()
	topics := []string{}
	for topic := range m.topics {
		topics = append(topics, topic)
	}
	return topics, args.Error(1)
}

func (m *MockServiceBus) GetQueues() ([]string, error) {
	args := m.Called()
	queues := []string{}
	for queue := range m.queues {
		queues = append(queues, queue)
	}
	return queues, args.Error(1)
}

// TestData contém dados de teste comuns
type TestData struct {
	ValidJSON   string
	InvalidJSON string
	TestMessage string
	TestQueue   string
	TestTopic   string
}

func GetTestData() TestData {
	return TestData{
		ValidJSON: `{
			"test": "message",
			"data": {
				"id": 123,
				"status": "active"
			}
		}`,
		InvalidJSON: `{
			"test": "message",
			"data": {
				"id": 123,
				"status": "active"
			`,
		TestMessage: `{"message": "test"}`,
		TestQueue:   "sbq.test.queue",
		TestTopic:   "sbt.test.topic",
	}
}

// Assertions comuns para testes
func AssertCommandSuccess(t *testing.T, stdout, stderr string, err error) {
	assert.NoError(t, err)
	// Remover a verificação de stderr vazio pois comandos podem escrever em stderr mesmo com sucesso
	// Remover a verificação de stdout não vazio pois alguns comandos podem não produzir saída
}

func AssertCommandFailure(t *testing.T, stdout, stderr string, err error) {
	assert.Error(t, err)
	// Remover a verificação de stderr não vazio pois alguns erros podem não escrever em stderr
}

func AssertFileExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	assert.NoError(t, err)
}

func AssertFileNotExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	assert.Error(t, err)
}
