package commands

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Parar ambiente",
	Long:  `Para o ambiente de testes Orion e limpa recursos.`,
	RunE:  runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)

	_, _ = blue.Println("🛑 Parando ambiente de testes Orion...")
	fmt.Println()

	// Parar containers
	if err := stopContainers(); err != nil {
		return err
	}

	// Remover containers órfãos
	if err := removeOrphanContainers(); err != nil {
		return err
	}

	// Limpar recursos se solicitado
	cleanFlag, _ := cmd.Flags().GetBool("clean")
	if cleanFlag {
		if err := cleanResources(); err != nil {
			return err
		}
	}

	// Mostrar informações finais
	showFinalInfoStop()

	return nil
}

func stopContainers() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Parando containers...")
	cmd := exec.Command("docker-compose", "down")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao parar containers: %w", err)
	}
	_, _ = green.Println("Containers parados")
	return nil
}

func removeOrphanContainers() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Removendo containers órfãos...")
	cmd := exec.Command("docker", "container", "prune", "-f")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run() // Ignorar erros aqui
	_, _ = green.Println("Containers órfãos removidos")
	return nil
}

func cleanResources() error {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	_, _ = yellow.Println("Limpando recursos (volumes e imagens)...")

	// Parar e remover volumes
	cmd := exec.Command("docker-compose", "down", "-v", "--rmi", "all")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run()

	// Limpar sistema Docker
	cmd = exec.Command("docker", "system", "prune", "-f")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run()

	_, _ = green.Println("Recursos limpos")
	return nil
}

func showFinalInfoStop() {
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	fmt.Println()
	_, _ = green.Println("🛑 Ambiente parado com sucesso!")
	fmt.Println()
	_, _ = blue.Println("📋 Para reiniciar o ambiente:")
	fmt.Println("  - orion-dev start")
	fmt.Println()
	_, _ = blue.Println("🧹 Para limpeza completa:")
	fmt.Println("  - orion-dev stop --clean")
	fmt.Println()
}

func init() {
	stopCmd.Flags().BoolP("clean", "c", false, "Limpar recursos (volumes e imagens)")
}
