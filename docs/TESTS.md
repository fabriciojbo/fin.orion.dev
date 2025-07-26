# 🧪 Testes do Fin.Orion.Dev

Este diretório contém testes mocados para todas as funcionalidades e comandos do projeto **Fin.Orion.Dev**.

## 📋 Estrutura dos Testes

```
Fin.Orion.Dev/
├── tests/                 # Testes Go
│   ├── main_test.go       # Utilitários e mocks comuns
│   ├── setup_test.go      # Testes do comando setup
│   ├── messages_test.go   # Testes dos comandos de mensagens
│   ├── environment_test.go # Testes dos comandos de ambiente
│   └── utils_test.go      # Testes de utilitários
├── scripts/               # Scripts de automação
│   └── run_tests.sh       # Script para executar testes
├── .github/workflows/     # CI/CD
│   └── tests.yml          # Workflow de testes
├── Makefile               # Comandos make para testes
└── README.tests.md        # Este arquivo
```

## 🚀 Como Executar os Testes

### Pré-requisitos

```bash
# Instalar dependências de teste
cd tests
go mod tidy
cd ..
```

### Executar Todos os Testes

```bash
# Usando o script
./scripts/run_tests.sh

# Usando make
make test

# Executar com cobertura
make test-coverage

# Executar com verbose
make test-verbose
```

### Executar Testes Específicos

```bash
# Usando o script
./scripts/run_tests.sh -s  # Testes de setup
./scripts/run_tests.sh -m  # Testes de mensagens
./scripts/run_tests.sh -e  # Testes de ambiente

# Usando make
make test-setup
make test-messages
make test-environment

# Teste específico
cd tests && go test -v -run TestSetupCommand && cd ..
```

## 📊 Cobertura de Testes

### Comandos Testados

#### 🔧 Comandos de Setup
- ✅ `setup` - Configuração inicial do ambiente
- ✅ Verificação de dependências
- ✅ Geração de certificados
- ✅ Criação de arquivos `.env` e `local.settings.json`
- ✅ Verificação da estrutura do projeto

#### 📨 Comandos de Mensagens
- ✅ `push-message` - Enviar mensagem para fila
- ✅ `check-messages` - Verificar mensagens do Service Bus
- ✅ `check-queue` - Verificar mensagens de fila específica
- ✅ `check-topic` - Verificar mensagens do tópico
- ✅ `list` - Listar recursos disponíveis
- ✅ `validate-json` - Validar arquivo JSON
- ✅ `format-json` - Formatar arquivo JSON
- ✅ `show-json` - Mostrar JSON formatado
- ✅ `proxy` - Iniciar proxy Service Bus

#### 🐳 Comandos de Ambiente
- ✅ `start` - Iniciar ambiente completo
- ✅ `stop` - Parar ambiente
- ✅ `status` - Verificar status dos containers
- ✅ `restart` - Reiniciar ambiente
- ✅ `logs` - Ver logs dos containers
- ✅ `clean` - Limpar recursos
- ✅ `build` - Reconstruir containers
- ✅ `shell` - Acessar shell dos containers
- ✅ `dev` - Modo desenvolvimento
- ✅ `quick-test` - Teste rápido
- ✅ `monitor` - Monitorar ambiente
- ✅ `health` - Verificar saúde dos serviços
- ✅ `rebuild-api` - Reconstruir Orion API
- ✅ `rebuild-functions` - Reconstruir Orion Functions
- ✅ `debug` - Debug do ambiente
- ✅ `debug-functions` - Debug do Orion Functions
- ✅ `clean-volumes` - Limpar volumes Docker
- ✅ `clean-images` - Limpar imagens Docker

## 🎯 Tipos de Teste

### 1. Testes Unitários
- Testam funcionalidades específicas isoladamente
- Usam mocks para simular dependências externas
- Verificam comportamento esperado com diferentes inputs

### 2. Testes de Integração
- Testam fluxos completos de comandos
- Simulam cenários reais de uso
- Verificam interação entre diferentes componentes

### 3. Testes de Cenários
- Testam casos de sucesso e falha
- Verificam validação de argumentos
- Testam tratamento de erros

## 🔧 Mocks Disponíveis

### MockFileSystem
```go
mockFS := NewMockFileSystem()
mockFS.On("WriteFile", filename, data, perm).Return(nil)
mockFS.On("ReadFile", filename).Return(content, nil)
mockFS.On("Stat", filename).Return(fileInfo, nil)
mockFS.On("MkdirAll", path, perm).Return(nil)
```

### MockDocker
```go
mockDocker := NewMockDocker()
mockDocker.On("Run", args).Return(nil)
mockDocker.SetContainerStatus("container", "running")
```

### MockServiceBus
```go
mockSB := NewMockServiceBus()
mockSB.On("SendMessage", queue, message).Return(nil)
mockSB.On("GetMessages", queue).Return(messages, nil)
mockSB.On("GetQueues").Return(queues, nil)
mockSB.On("GetTopics").Return(topics, nil)
```

## 📝 Exemplos de Uso

### Teste Simples
```go
func TestExample(t *testing.T) {
    TestSetup(t)

    t.Run("cenário de sucesso", func(t *testing.T) {
        // Arrange
        cmd := &cobra.Command{
            Use: "test",
            RunE: func(cmd *cobra.Command, args []string) error {
                return nil
            },
        }

        // Act
        stdout, stderr, err := executeCommand(t, cmd)

        // Assert
        AssertCommandSuccess(t, stdout, stderr, err)
    })
}
```

### Teste com Mocks
```go
func TestWithMocks(t *testing.T) {
    TestSetup(t)

    t.Run("teste com mock", func(t *testing.T) {
        // Arrange
        mockFS := NewMockFileSystem()
        mockFS.On("WriteFile", "test.txt", []byte("content"), os.FileMode(0644)).Return(nil)

        // Act & Assert
        assert.True(t, true)
    })
}
```

## 🚨 Troubleshooting

### Erro: "missing metadata for import"
```bash
# Solução: Executar go mod tidy
cd tests
go mod tidy
```

### Erro: "undefined: TestSetup"
```bash
# Solução: Verificar se o arquivo main_test.go está presente
ls -la tests/main_test.go
```

### Erro: "command not found"
```bash
# Solução: Instalar dependências
cd tests
go mod download
```

## 📈 Melhorias Futuras

- [ ] Adicionar testes de performance
- [ ] Implementar testes de stress
- [ ] Adicionar testes de segurança
- [ ] Implementar testes de regressão
- [ ] Adicionar testes de compatibilidade

## 🤝 Contribuindo

Para adicionar novos testes:

1. Crie um arquivo `*_test.go` seguindo a convenção
2. Use os mocks disponíveis em `main_test.go`
3. Siga o padrão Arrange-Act-Assert
4. Adicione documentação clara
5. Execute `go test` para verificar

## 📚 Recursos

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Testify Framework](https://github.com/stretchr/testify)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Go Testing Best Practices](https://golang.org/doc/tutorial/add-a-test)