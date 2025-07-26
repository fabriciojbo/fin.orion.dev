package commands

import (
	"fin.orion.dev/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "orion-dev",
	Short: "Orion Development Environment Tools",
	Long: `🔧 Ferramentas para o ambiente de desenvolvimento Orion Functions

Este CLI fornece comandos para gerenciar o ambiente de desenvolvimento,
enviar mensagens para filas e tópicos do Service Bus, e monitorar
o status dos serviços.`,
	Version: utils.GetVersionOrUnknown(),
	Run: func(cmd *cobra.Command, args []string) {
		// Se não houver argumentos, mostrar ajuda
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

// Execute adiciona todos os comandos filhos ao comando raiz e configura flags
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Adicionar subcomandos de ambiente
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(quickTestCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(rebuildApiCmd)
	rootCmd.AddCommand(rebuildFunctionsCmd)
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(debugFunctionsCmd)
	rootCmd.AddCommand(cleanVolumesCmd)
	rootCmd.AddCommand(cleanImagesCmd)

	// Adicionar subcomandos de mensagens
	rootCmd.AddCommand(checkMessagesCmd)
	rootCmd.AddCommand(pushMessageCmd)
	rootCmd.AddCommand(checkQueueCmd)
	rootCmd.AddCommand(checkTopicCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(testMessageCmd)
	rootCmd.AddCommand(sendQueueCmd)
	rootCmd.AddCommand(sendJsonCmd)
	rootCmd.AddCommand(listQueuesCmd)
	rootCmd.AddCommand(listMessagesCmd)
	rootCmd.AddCommand(validateJsonCmd)
	rootCmd.AddCommand(formatJsonCmd)
	rootCmd.AddCommand(showJsonCmd)
	rootCmd.AddCommand(proxyCmd)

	// Adicionar subcomandos de commitlint
	rootCmd.AddCommand(commitlintCmd)
	rootCmd.AddCommand(commitlintLastCmd)
	rootCmd.AddCommand(commitlintTypesCmd)
	rootCmd.AddCommand(commitlintFormatCmd)
	rootCmd.AddCommand(commitlintHookCmd)
}
