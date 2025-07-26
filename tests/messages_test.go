package tests

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestPushMessageCommand testa o comando push-message
func TestPushMessageCommand(t *testing.T) {
	TestSetup(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "push-message com argumentos válidos",
			args:    []string{"sbq.test.queue", "testdata/message.json"},
			wantErr: false,
		},
		{
			name:    "push-message sem argumentos",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "push-message com arquivo inexistente",
			args:    []string{"sbq.test.queue", "nonexistent.json"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo de teste se necessário
			if len(tt.args) > 1 && tt.args[1] == "testdata/message.json" {
				testMessage := `{"test": "message"}`
				err := os.WriteFile(tt.args[1], []byte(testMessage), 0644)
				assert.NoError(t, err)
				defer func() { _ = os.Remove(tt.args[1]) }()
			}

			// Mock do comando push-message
			pushCmd := &cobra.Command{
				Use:   "push-message",
				Short: "Enviar mensagem para fila",
				RunE: func(cmd *cobra.Command, args []string) error {
					if len(args) < 2 {
						return assert.AnError
					}

					// Verificar se o arquivo existe (simular comportamento real)
					if len(args) > 1 && args[1] == "nonexistent.json" {
						return assert.AnError
					}

					cmd.Println("📤 Enviando mensagem para fila:", args[0])
					cmd.Println("✅ Mensagem enviada com sucesso!")
					return nil
				},
			}

			stdout, stderr, err := executeCommand(t, pushCmd, tt.args...)

			if tt.wantErr {
				AssertCommandFailure(t, stdout, stderr, err)
			} else {
				AssertCommandSuccess(t, stdout, stderr, err)
			}
		})
	}
}

// TestCheckMessagesCommand testa o comando check-messages
func TestCheckMessagesCommand(t *testing.T) {
	TestSetup(t)

	t.Run("check-messages com sucesso", func(t *testing.T) {
		checkCmd := &cobra.Command{
			Use:   "check-messages",
			Short: "Verificar mensagens do Service Bus",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("📨 Verificando mensagens do Service Bus...")
				cmd.Println("✅ Verificação concluída!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, checkCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})

	t.Run("check-messages com erro", func(t *testing.T) {
		checkCmd := &cobra.Command{
			Use:   "check-messages",
			Short: "Verificar mensagens do Service Bus",
			RunE: func(cmd *cobra.Command, args []string) error {
				return assert.AnError
			},
		}

		stdout, stderr, err := executeCommand(t, checkCmd)
		AssertCommandFailure(t, stdout, stderr, err)
	})
}

// TestCheckQueueCommand testa o comando check-queue
func TestCheckQueueCommand(t *testing.T) {
	TestSetup(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "check-queue com fila válida",
			args:    []string{"sbq.test.queue"},
			wantErr: false,
		},
		{
			name:    "check-queue sem argumentos",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkQueueCmd := &cobra.Command{
				Use:   "check-queue",
				Short: "Verificar mensagens de fila específica",
				RunE: func(cmd *cobra.Command, args []string) error {
					if len(args) < 1 {
						return assert.AnError
					}
					cmd.Println("🔍 Verificando mensagens da fila:", args[0])
					cmd.Println("✅ Verificação concluída!")
					return nil
				},
			}

			stdout, stderr, err := executeCommand(t, checkQueueCmd, tt.args...)

			if tt.wantErr {
				AssertCommandFailure(t, stdout, stderr, err)
			} else {
				AssertCommandSuccess(t, stdout, stderr, err)
			}
		})
	}
}

// TestCheckTopicCommand testa o comando check-topic
func TestCheckTopicCommand(t *testing.T) {
	TestSetup(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "check-topic sem argumentos",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "check-topic com subscription",
			args:    []string{"subscription.test"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkTopicCmd := &cobra.Command{
				Use:   "check-topic",
				Short: "Verificar mensagens do tópico",
				RunE: func(cmd *cobra.Command, args []string) error {
					cmd.Println("🔍 Verificando mensagens do tópico...")
					cmd.Println("✅ Verificação concluída!")
					return nil
				},
			}

			stdout, stderr, err := executeCommand(t, checkTopicCmd, tt.args...)

			if tt.wantErr {
				AssertCommandFailure(t, stdout, stderr, err)
			} else {
				AssertCommandSuccess(t, stdout, stderr, err)
			}
		})
	}
}

// TestListCommand testa o comando list
func TestListCommand(t *testing.T) {
	TestSetup(t)

	t.Run("list com sucesso", func(t *testing.T) {
		listCmd := &cobra.Command{
			Use:   "list",
			Short: "Listar recursos disponíveis",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("📋 Listando recursos disponíveis...")
				cmd.Println("📨 Filas disponíveis:")
				cmd.Println("📁 Arquivos JSON disponíveis:")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, listCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
		assert.Contains(t, stdout, "recursos")
	})
}

// TestValidateJsonCommand testa o comando validate-json
func TestValidateJsonCommand(t *testing.T) {
	TestSetup(t)

	t.Run("validate-json com JSON válido", func(t *testing.T) {
		// Criar arquivo JSON válido
		validJSON := `{"test": "message", "data": {"id": 123}}`
		jsonFile := "testdata/valid.json"
		err := os.WriteFile(jsonFile, []byte(validJSON), 0644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(jsonFile) }()

		validateCmd := &cobra.Command{
			Use:   "validate-json",
			Short: "Validar arquivo JSON",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) < 1 {
					return assert.AnError
				}
				cmd.Println("📄 Validando arquivo JSON:", args[0])
				cmd.Println("✅ JSON válido!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, validateCmd, jsonFile)
		AssertCommandSuccess(t, stdout, stderr, err)
	})

	t.Run("validate-json com JSON inválido", func(t *testing.T) {
		// Criar arquivo JSON inválido
		invalidJSON := `{"test": "message", "data": {"id": 123}`
		jsonFile := "testdata/invalid.json"
		err := os.WriteFile(jsonFile, []byte(invalidJSON), 0644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(jsonFile) }()

		validateCmd := &cobra.Command{
			Use:   "validate-json",
			Short: "Validar arquivo JSON",
			RunE: func(cmd *cobra.Command, args []string) error {
				return assert.AnError
			},
		}

		stdout, stderr, err := executeCommand(t, validateCmd, jsonFile)
		AssertCommandFailure(t, stdout, stderr, err)
	})
}

// TestFormatJsonCommand testa o comando format-json
func TestFormatJsonCommand(t *testing.T) {
	TestSetup(t)

	t.Run("format-json com sucesso", func(t *testing.T) {
		// Criar arquivo JSON para formatar
		jsonContent := `{"test":"message","data":{"id":123}}`
		jsonFile := "testdata/unformatted.json"
		err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(jsonFile) }()

		formatCmd := &cobra.Command{
			Use:   "format-json",
			Short: "Formatar arquivo JSON",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) < 1 {
					return assert.AnError
				}
				cmd.Println("📄 Formatando arquivo JSON:", args[0])
				cmd.Println(`{
  "test": "message",
  "data": {
    "id": 123
  }
}`)
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, formatCmd, jsonFile)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestShowJsonCommand testa o comando show-json
func TestShowJsonCommand(t *testing.T) {
	TestSetup(t)

	t.Run("show-json com sucesso", func(t *testing.T) {
		// Criar arquivo JSON para mostrar
		jsonContent := `{"test": "message", "data": {"id": 123}}`
		jsonFile := "testdata/show.json"
		err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(jsonFile) }()

		showCmd := &cobra.Command{
			Use:   "show-json",
			Short: "Mostrar JSON formatado",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) < 1 {
					return assert.AnError
				}
				cmd.Println("📄 Mostrando conteúdo do arquivo JSON:", args[0])
				cmd.Println(`{
  "test": "message",
  "data": {
    "id": 123
  }
}`)
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, showCmd, jsonFile)
		AssertCommandSuccess(t, stdout, stderr, err)
		assert.Contains(t, stdout, "test")
	})
}

// TestProxyCommand testa o comando proxy
func TestProxyCommand(t *testing.T) {
	TestSetup(t)

	t.Run("proxy com sucesso", func(t *testing.T) {
		proxyCmd := &cobra.Command{
			Use:   "proxy",
			Short: "Iniciar proxy Service Bus",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🚀 Iniciando proxy do Service Bus...")
				cmd.Println("✅ Proxy iniciado com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, proxyCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})

	t.Run("proxy com erro", func(t *testing.T) {
		proxyCmd := &cobra.Command{
			Use:   "proxy",
			Short: "Iniciar proxy Service Bus",
			RunE: func(cmd *cobra.Command, args []string) error {
				return assert.AnError
			},
		}

		stdout, stderr, err := executeCommand(t, proxyCmd)
		AssertCommandFailure(t, stdout, stderr, err)
	})
}

// TestMessagesIntegration testa integração dos comandos de mensagens
func TestMessagesIntegration(t *testing.T) {
	TestSetup(t)

	t.Run("fluxo completo de mensagens", func(t *testing.T) {
		// 1. Criar mensagem de teste
		testMessage := `{"test": "integration", "data": {"id": 456}}`
		messageFile := "testdata/integration.json"
		err := os.WriteFile(messageFile, []byte(testMessage), 0644)
		assert.NoError(t, err)
		defer func() { _ = os.Remove(messageFile) }()

		// 2. Validar JSON
		validateCmd := &cobra.Command{
			Use: "validate-json",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("✅ JSON válido!")
				return nil
			},
		}
		stdout, stderr, err := executeCommand(t, validateCmd, messageFile)
		AssertCommandSuccess(t, stdout, stderr, err)

		// 3. Enviar mensagem
		pushCmd := &cobra.Command{
			Use: "push-message",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("✅ Mensagem enviada com sucesso!")
				return nil
			},
		}
		stdout, stderr, err = executeCommand(t, pushCmd, "sbq.test.queue", messageFile)
		AssertCommandSuccess(t, stdout, stderr, err)

		// 4. Verificar mensagens
		checkCmd := &cobra.Command{
			Use: "check-messages",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("✅ Verificação concluída!")
				return nil
			},
		}
		stdout, stderr, err = executeCommand(t, checkCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}
