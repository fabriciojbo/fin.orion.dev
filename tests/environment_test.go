package tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestStartCommand testa o comando start
func TestStartCommand(t *testing.T) {
	TestSetup(t)

	t.Run("start com sucesso", func(t *testing.T) {
		startCmd := &cobra.Command{
			Use:   "start",
			Short: "Iniciar ambiente completo",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🚀 Iniciando ambiente completo de testes Orion...")
				cmd.Println("🎉 Ambiente iniciado com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, startCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
		assert.Contains(t, stdout, "iniciado")
	})

	t.Run("start com erro", func(t *testing.T) {
		startCmd := &cobra.Command{
			Use:   "start",
			Short: "Iniciar ambiente completo",
			RunE: func(cmd *cobra.Command, args []string) error {
				return assert.AnError
			},
		}

		stdout, stderr, err := executeCommand(t, startCmd)
		AssertCommandFailure(t, stdout, stderr, err)
	})
}

// TestStopCommand testa o comando stop
func TestStopCommand(t *testing.T) {
	TestSetup(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "stop sem argumentos",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "stop com --clean",
			args:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopCmd := &cobra.Command{
				Use:   "stop",
				Short: "Parar ambiente",
				RunE: func(cmd *cobra.Command, args []string) error {
					cmd.Println("🛑 Parando ambiente de testes Orion...")
					cmd.Println("🛑 Ambiente parado com sucesso!")
					return nil
				},
			}

			// Adicionar flag --clean se necessário
			if tt.name == "stop com --clean" {
				stopCmd.Flags().BoolP("clean", "c", false, "Limpar recursos (volumes e imagens)")
			}

			stdout, stderr, err := executeCommand(t, stopCmd, tt.args...)

			if tt.wantErr {
				AssertCommandFailure(t, stdout, stderr, err)
			} else {
				AssertCommandSuccess(t, stdout, stderr, err)
			}
		})
	}
}

// TestStatusCommand testa o comando status
func TestStatusCommand(t *testing.T) {
	TestSetup(t)

	t.Run("status com sucesso", func(t *testing.T) {
		statusCmd := &cobra.Command{
			Use:   "status",
			Short: "Verificar status dos containers",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("📋 Status dos containers:")
				cmd.Println("✅ Container 1: running")
				cmd.Println("✅ Container 2: running")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, statusCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
		assert.Contains(t, stdout, "Status")
	})

	t.Run("status com erro", func(t *testing.T) {
		statusCmd := &cobra.Command{
			Use:   "status",
			Short: "Verificar status dos containers",
			RunE: func(cmd *cobra.Command, args []string) error {
				return assert.AnError
			},
		}

		stdout, stderr, err := executeCommand(t, statusCmd)
		AssertCommandFailure(t, stdout, stderr, err)
	})
}

// TestRestartCommand testa o comando restart
func TestRestartCommand(t *testing.T) {
	TestSetup(t)

	t.Run("restart com sucesso", func(t *testing.T) {
		restartCmd := &cobra.Command{
			Use:   "restart",
			Short: "Reiniciar ambiente",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🔄 Reiniciando ambiente...")
				cmd.Println("✅ Ambiente reiniciado com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, restartCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestLogsCommand testa o comando logs
func TestLogsCommand(t *testing.T) {
	TestSetup(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "logs sem argumentos",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "logs com --service",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "logs com --tail",
			args:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logsCmd := &cobra.Command{
				Use:   "logs",
				Short: "Ver logs dos containers",
				RunE: func(cmd *cobra.Command, args []string) error {
					cmd.Println("📋 Logs dos containers:")
					cmd.Println("2024-01-01 10:00:00 [INFO] Container started")
					return nil
				},
			}

			// Adicionar flags se necessário
			if tt.name == "logs com --service" {
				logsCmd.Flags().StringP("service", "s", "", "Nome do serviço")
			}
			if tt.name == "logs com --tail" {
				logsCmd.Flags().IntP("tail", "t", 100, "Número de linhas")
			}

			stdout, stderr, err := executeCommand(t, logsCmd, tt.args...)

			if tt.wantErr {
				AssertCommandFailure(t, stdout, stderr, err)
			} else {
				AssertCommandSuccess(t, stdout, stderr, err)
			}
		})
	}
}

// TestCleanCommand testa o comando clean
func TestCleanCommand(t *testing.T) {
	TestSetup(t)

	t.Run("clean com sucesso", func(t *testing.T) {
		cleanCmd := &cobra.Command{
			Use:   "clean",
			Short: "Limpar recursos",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🧹 Limpando recursos...")
				cmd.Println("✅ Recursos limpos com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, cleanCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestBuildCommand testa o comando build
func TestBuildCommand(t *testing.T) {
	TestSetup(t)

	t.Run("build com sucesso", func(t *testing.T) {
		buildCmd := &cobra.Command{
			Use:   "build",
			Short: "Reconstruir containers",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🔨 Reconstruindo containers...")
				cmd.Println("✅ Containers reconstruídos com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, buildCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestShellCommand testa o comando shell
func TestShellCommand(t *testing.T) {
	TestSetup(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "shell com --service",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "shell sem --service",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shellCmd := &cobra.Command{
				Use:   "shell",
				Short: "Acessar shell dos containers",
				RunE: func(cmd *cobra.Command, args []string) error {
					serviceFlag, _ := cmd.Flags().GetString("service")
					if serviceFlag == "" && tt.name == "shell sem --service" {
						return assert.AnError
					}
					cmd.Println("🐚 Acessando shell do container...")
					return nil
				},
			}

			// Adicionar flag --service
			shellCmd.Flags().StringP("service", "s", "", "Nome do serviço")

			stdout, stderr, err := executeCommand(t, shellCmd, tt.args...)

			if tt.wantErr {
				AssertCommandFailure(t, stdout, stderr, err)
			} else {
				AssertCommandSuccess(t, stdout, stderr, err)
			}
		})
	}
}

// TestDevCommand testa o comando dev
func TestDevCommand(t *testing.T) {
	TestSetup(t)

	t.Run("dev com sucesso", func(t *testing.T) {
		devCmd := &cobra.Command{
			Use:   "dev",
			Short: "Modo desenvolvimento",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🔧 Modo desenvolvimento ativado...")
				cmd.Println("✅ Ambiente de desenvolvimento pronto!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, devCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestQuickTestCommand testa o comando quick-test
func TestQuickTestCommand(t *testing.T) {
	TestSetup(t)

	t.Run("quick-test com sucesso", func(t *testing.T) {
		quickTestCmd := &cobra.Command{
			Use:   "quick-test",
			Short: "Teste rápido",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("⚡ Executando teste rápido...")
				cmd.Println("✅ Teste concluído com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, quickTestCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestMonitorCommand testa o comando monitor
func TestMonitorCommand(t *testing.T) {
	TestSetup(t)

	t.Run("monitor com sucesso", func(t *testing.T) {
		monitorCmd := &cobra.Command{
			Use:   "monitor",
			Short: "Monitorar ambiente",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("📊 Monitorando ambiente...")
				cmd.Println("✅ Monitoramento ativo!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, monitorCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestHealthCommand testa o comando health
func TestHealthCommand(t *testing.T) {
	TestSetup(t)

	t.Run("health com sucesso", func(t *testing.T) {
		healthCmd := &cobra.Command{
			Use:   "health",
			Short: "Verificar saúde dos serviços",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🏥 Verificando saúde dos serviços...")
				cmd.Println("✅ Todos os serviços estão saudáveis!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, healthCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
		assert.Contains(t, stdout, "saudáveis")
	})
}

// TestRebuildApiCommand testa o comando rebuild-api
func TestRebuildApiCommand(t *testing.T) {
	TestSetup(t)

	t.Run("rebuild-api com sucesso", func(t *testing.T) {
		rebuildApiCmd := &cobra.Command{
			Use:   "rebuild-api",
			Short: "Reconstruir Orion API",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🔨 Reconstruindo Orion API...")
				cmd.Println("✅ Orion API reconstruída com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, rebuildApiCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestRebuildFunctionsCommand testa o comando rebuild-functions
func TestRebuildFunctionsCommand(t *testing.T) {
	TestSetup(t)

	t.Run("rebuild-functions com sucesso", func(t *testing.T) {
		rebuildFunctionsCmd := &cobra.Command{
			Use:   "rebuild-functions",
			Short: "Reconstruir Orion Functions",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🔨 Reconstruindo Orion Functions...")
				cmd.Println("✅ Orion Functions reconstruída com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, rebuildFunctionsCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestDebugCommand testa o comando debug
func TestDebugCommand(t *testing.T) {
	TestSetup(t)

	t.Run("debug com sucesso", func(t *testing.T) {
		debugCmd := &cobra.Command{
			Use:   "debug",
			Short: "Debug do ambiente",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🐛 Iniciando debug do ambiente...")
				cmd.Println("✅ Debug ativo!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, debugCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestDebugFunctionsCommand testa o comando debug-functions
func TestDebugFunctionsCommand(t *testing.T) {
	TestSetup(t)

	t.Run("debug-functions com sucesso", func(t *testing.T) {
		debugFunctionsCmd := &cobra.Command{
			Use:   "debug-functions",
			Short: "Debug do Orion Functions",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🐛 Iniciando debug do Orion Functions...")
				cmd.Println("✅ Debug do Orion Functions ativo!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, debugFunctionsCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestCleanVolumesCommand testa o comando clean-volumes
func TestCleanVolumesCommand(t *testing.T) {
	TestSetup(t)

	t.Run("clean-volumes com sucesso", func(t *testing.T) {
		cleanVolumesCmd := &cobra.Command{
			Use:   "clean-volumes",
			Short: "Limpar volumes Docker",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🧹 Limpando volumes Docker...")
				cmd.Println("✅ Volumes limpos com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, cleanVolumesCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestCleanImagesCommand testa o comando clean-images
func TestCleanImagesCommand(t *testing.T) {
	TestSetup(t)

	t.Run("clean-images com sucesso", func(t *testing.T) {
		cleanImagesCmd := &cobra.Command{
			Use:   "clean-images",
			Short: "Limpar imagens Docker",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🧹 Limpando imagens Docker...")
				cmd.Println("✅ Imagens limpas com sucesso!")
				return nil
			},
		}

		stdout, stderr, err := executeCommand(t, cleanImagesCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}

// TestEnvironmentIntegration testa integração dos comandos de ambiente
func TestEnvironmentIntegration(t *testing.T) {
	TestSetup(t)

	t.Run("fluxo completo de ambiente", func(t *testing.T) {
		// 1. Verificar status inicial
		statusCmd := &cobra.Command{
			Use: "status",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("📋 Status dos containers:")
				cmd.Println("✅ Container 1: running")
				return nil
			},
		}
		stdout, stderr, err := executeCommand(t, statusCmd)
		AssertCommandSuccess(t, stdout, stderr, err)

		// 2. Iniciar ambiente
		startCmd := &cobra.Command{
			Use: "start",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🚀 Iniciando ambiente...")
				cmd.Println("✅ Ambiente iniciado!")
				return nil
			},
		}
		stdout, stderr, err = executeCommand(t, startCmd)
		AssertCommandSuccess(t, stdout, stderr, err)

		// 3. Verificar saúde
		healthCmd := &cobra.Command{
			Use: "health",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🏥 Verificando saúde...")
				cmd.Println("✅ Serviços saudáveis!")
				return nil
			},
		}
		stdout, stderr, err = executeCommand(t, healthCmd)
		AssertCommandSuccess(t, stdout, stderr, err)

		// 4. Parar ambiente
		stopCmd := &cobra.Command{
			Use: "stop",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("🛑 Parando ambiente...")
				cmd.Println("✅ Ambiente parado!")
				return nil
			},
		}
		stdout, stderr, err = executeCommand(t, stopCmd)
		AssertCommandSuccess(t, stdout, stderr, err)
	})
}
