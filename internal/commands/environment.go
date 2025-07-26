package commands

import (
	"fmt"
	"os/exec"

	"fin.orion.dev/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Comando para reiniciar ambiente
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Reiniciar ambiente",
	Long:  `Para e reinicia o ambiente de testes Orion.`,
	RunE:  runRestart,
}

// Comando para ver logs
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Ver logs dos containers",
	Long:  `Mostra os logs dos containers em tempo real.`,
	RunE:  runLogs,
}

// Comando para limpar ambiente
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Limpar ambiente completamente",
	Long:  `Para o ambiente e remove todos os containers, volumes e imagens.`,
	RunE:  runClean,
}

// Comando para reconstruir containers
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Reconstruir containers",
	Long:  `Reconstrói todos os containers sem cache.`,
	RunE:  runBuild,
}

// Comando para acessar shell
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Acessar shell do Orion Functions",
	Long:  `Abre um shell interativo no container do Orion Functions.`,
	RunE:  runShell,
}

// Comando para desenvolvimento
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Iniciar ambiente para desenvolvimento",
	Long:  `Inicia o ambiente completo para desenvolvimento.`,
	RunE:  runDev,
}

// Comando para teste rápido
var quickTestCmd = &cobra.Command{
	Use:   "quick-test",
	Short: "Teste rápido das functions",
	Long:  `Executa testes rápidos para verificar se as functions estão funcionando.`,
	RunE:  runQuickTest,
}

// Comando para monitorar recursos
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitorar recursos do sistema",
	Long:  `Mostra estatísticas de uso de recursos dos containers.`,
	RunE:  runMonitor,
}

// Comando para verificar saúde
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Verificar saúde dos serviços",
	Long:  `Verifica a saúde de todos os serviços do ambiente.`,
	RunE:  runHealth,
}

// Comando para reconstruir API
var rebuildApiCmd = &cobra.Command{
	Use:   "rebuild-api",
	Short: "Reconstruir apenas o Orion API",
	Long:  `Reconstrói apenas o container do Orion API.`,
	RunE:  runRebuildApi,
}

// Comando para reconstruir Functions
var rebuildFunctionsCmd = &cobra.Command{
	Use:   "rebuild-functions",
	Short: "Reconstruir apenas o Orion Functions",
	Long:  `Reconstrói apenas o container do Orion Functions.`,
	RunE:  runRebuildFunctions,
}

// Comando para debug
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Modo debug - logs detalhados",
	Long:  `Mostra logs detalhados de todos os containers.`,
	RunE:  runDebug,
}

// Comando para debug específico do Functions
var debugFunctionsCmd = &cobra.Command{
	Use:   "debug-functions",
	Short: "Debug específico do Orion Functions",
	Long:  `Mostra logs detalhados apenas do Orion Functions.`,
	RunE:  runDebugFunctions,
}

// Comando para limpar volumes
var cleanVolumesCmd = &cobra.Command{
	Use:   "clean-volumes",
	Short: "Limpar apenas volumes",
	Long:  `Remove apenas os volumes do ambiente.`,
	RunE:  runCleanVolumes,
}

// Comando para limpar imagens
var cleanImagesCmd = &cobra.Command{
	Use:   "clean-images",
	Short: "Limpar imagens não utilizadas",
	Long:  `Remove imagens Docker não utilizadas.`,
	RunE:  runCleanImages,
}

func runRestart(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🔄 Reiniciando ambiente...")
	fmt.Println()

	// Parar ambiente
	if err := runStop(cmd, args); err != nil {
		return err
	}

	fmt.Println()

	// Iniciar ambiente
	if err := runStart(cmd, args); err != nil {
		return err
	}

	_, _ = green.Println("✅ Ambiente reiniciado com sucesso!")
	return nil
}

func runLogs(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("📝 Mostrando logs dos containers...")

	cmdExec := exec.Command("docker-compose", "logs", "-f")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	return cmdExec.Run()
}

func runClean(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	_, _ = blue.Println("🧹 Limpando ambiente completamente...")
	_, _ = yellow.Println("⚠️  Esta operação irá remover todos os containers, volumes e imagens!")

	// Parar e remover tudo
	cmdExec := exec.Command("docker-compose", "down", "-v", "--rmi", "all")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao limpar containers: %w", err)
	}

	// Limpar sistema Docker
	cmdExec = exec.Command("docker", "system", "prune", "-f")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	_ = cmdExec.Run()

	_, _ = green.Println("✅ Ambiente limpo completamente!")
	return nil
}

func runBuild(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🔨 Reconstruindo containers...")

	cmdExec := exec.Command("docker-compose", "build", "--no-cache")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao reconstruir containers: %w", err)
	}

	_, _ = green.Println("✅ Containers reconstruídos!")
	return nil
}

func runShell(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("🐚 Acessando shell do Orion Functions...")

	cmdExec := exec.Command("docker-compose", "exec", "orion-functions", "/bin/sh")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	return cmdExec.Run()
}

func runDev(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("💡 Iniciando ambiente para desenvolvimento...")

	if err := runStart(cmd, args); err != nil {
		return err
	}

	fmt.Println()
	_, _ = green.Println("💡 Ambiente pronto para desenvolvimento!")
	_, _ = blue.Println("📋 URLs disponíveis:")
	_, _ = blue.Println("  - Orion Functions: http://localhost:7071")
	_, _ = blue.Println("  - Orion API: http://localhost:3333")

	return nil
}

func runQuickTest(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Println("⚡ Executando teste rápido das functions...")

	// Testar Orion Functions QR Code COB
	if utils.CheckHTTPEndpoint("http://localhost:7071/cob/test-id") {
		_, _ = green.Println("✅ Orion Functions QR Code COB funcionando")
	} else {
		_, _ = red.Println("❌ Orion Functions QR Code COB falhou")
	}

	// Testar Orion Functions QR Code COBV
	if utils.CheckHTTPEndpoint("http://localhost:7071/cobv/test-id") {
		_, _ = green.Println("✅ Orion Functions QR Code COBV funcionando")
	} else {
		_, _ = red.Println("❌ Orion Functions QR Code COBV falhou")
	}

	// Testar Orion API
	if utils.CheckHTTPEndpoint("http://localhost:3333") {
		_, _ = green.Println("✅ Orion API funcionando")
	} else {
		_, _ = red.Println("❌ Orion API falhou")
	}

	return nil
}

func runMonitor(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("📊 Monitorando recursos do sistema...")

	cmdExec := exec.Command("docker", "stats", "--no-stream")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	return cmdExec.Run()
}

func runHealth(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Println("🏥 Verificando saúde dos serviços...")

	// Azure Service Bus
	if utils.CheckHTTPEndpoint("http://localhost:5300/health") {
		_, _ = green.Println("✅ Azure Service Bus: OK")
	} else {
		_, _ = red.Println("❌ Azure Service Bus: ERRO")
	}

	// Azure Storage
	if utils.CheckHTTPEndpoint("http://localhost:10000") {
		_, _ = green.Println("✅ Azure Storage: OK")
	} else {
		_, _ = red.Println("❌ Azure Storage: ERRO")
	}

	// PostgreSQL
	cmdExec := exec.Command("docker-compose", "exec", "-T", "postgres", "pg_isready", "-U", "postgres")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if cmdExec.Run() == nil {
		_, _ = green.Println("✅ PostgreSQL: OK")
	} else {
		_, _ = red.Println("❌ PostgreSQL: ERRO")
	}

	// Orion API
	if utils.CheckHTTPEndpoint("http://localhost:3333") {
		_, _ = green.Println("✅ Orion API: OK")
	} else {
		_, _ = red.Println("❌ Orion API: ERRO")
	}

	// Orion Functions
	if utils.CheckHTTPEndpoint("http://localhost:7071") {
		_, _ = green.Println("✅ Orion Functions: OK")
	} else {
		_, _ = red.Println("❌ Orion Functions: ERRO")
	}

	return nil
}

func runRebuildApi(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🔨 Reconstruindo Orion API...")

	cmdExec := exec.Command("docker-compose", "build", "orion-api", "--no-cache")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao reconstruir Orion API: %w", err)
	}

	cmdExec = exec.Command("docker-compose", "up", "-d", "orion-api")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao reiniciar Orion API: %w", err)
	}

	_, _ = green.Println("✅ Orion API reconstruído!")
	return nil
}

func runRebuildFunctions(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🔨 Reconstruindo Orion Functions...")

	cmdExec := exec.Command("docker-compose", "build", "orion-functions", "--no-cache")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao reconstruir Orion Functions: %w", err)
	}

	cmdExec = exec.Command("docker-compose", "up", "-d", "orion-functions")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao reiniciar Orion Functions: %w", err)
	}

	_, _ = green.Println("✅ Orion Functions reconstruído!")
	return nil
}

func runDebug(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("🐛 Modo debug ativado - logs detalhados...")

	cmdExec := exec.Command("docker-compose", "logs", "-f", "--tail=100")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	return cmdExec.Run()
}

func runDebugFunctions(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("🐛 Debug do Orion Functions...")

	cmdExec := exec.Command("docker-compose", "logs", "-f", "--tail=50", "orion-functions")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	return cmdExec.Run()
}

func runCleanVolumes(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🧹 Limpando volumes...")

	cmdExec := exec.Command("docker-compose", "down", "-v")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao limpar volumes: %w", err)
	}

	_, _ = green.Println("✅ Volumes limpos!")
	return nil
}

func runCleanImages(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🧹 Limpando imagens não utilizadas...")

	cmdExec := exec.Command("docker", "image", "prune", "-f")
	cmdExec.Stdout = nil
	cmdExec.Stderr = nil
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("erro ao limpar imagens: %w", err)
	}

	_, _ = green.Println("✅ Imagens limpas!")
	return nil
}
