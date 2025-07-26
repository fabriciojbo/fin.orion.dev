package commands

import (
	"fmt"
	"os/exec"
	"time"

	"fin.orion.dev/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Iniciar ambiente completo",
	Long:  `Inicia o ambiente completo de testes Orion com todos os serviços.`,
	RunE:  runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)

	_, _ = blue.Println("🚀 Iniciando ambiente completo de testes Orion...")
	fmt.Println()

	// Verificar Docker
	if err := checkDocker(); err != nil {
		return err
	}

	// Parar containers existentes
	if err := stopExistingContainers(); err != nil {
		return err
	}

	// Limpar containers órfãos
	if err := cleanOrphanContainers(); err != nil {
		return err
	}

	// Construir e iniciar serviços
	if err := startServices(); err != nil {
		return err
	}

	// Aguardar inicialização
	waitForServices()

	// Verificar status dos containers
	checkContainerStatusStart()

	// Verificar conectividade dos serviços
	checkServiceConnectivity()

	// Mostrar informações finais
	showFinalInfoStart()

	return nil
}

func checkDocker() error {
	red := color.New(color.FgRed)
	blue := color.New(color.FgBlue)

	_, _ = blue.Println("Verificando Docker...")
	if err := exec.Command("docker", "info").Run(); err != nil {
		_, _ = red.Println("Docker não está rodando. Por favor, inicie o Docker e tente novamente.")
		return fmt.Errorf("docker não está rodando")
	}
	return nil
}

func stopExistingContainers() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Parando containers existentes...")
	cmd := exec.Command("docker-compose", "down", "--remove-orphans")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run() // Ignorar erros aqui
	_, _ = green.Println("Containers existentes parados")
	return nil
}

func cleanOrphanContainers() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Limpando containers órfãos...")
	cmd := exec.Command("docker", "container", "prune", "-f")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run() // Ignorar erros aqui
	_, _ = green.Println("Containers órfãos removidos")
	return nil
}

func startServices() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Construindo e iniciando serviços...")
	cmd := exec.Command("docker-compose", "up", "--build", "-d")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao iniciar serviços: %w", err)
	}
	_, _ = green.Println("Serviços iniciados")
	return nil
}

func waitForServices() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Aguardando serviços iniciarem...")
	time.Sleep(30 * time.Second)
	_, _ = green.Println("Tempo de espera concluído")
}

func checkContainerStatusStart() {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("Verificando status dos containers:")
	cmd := exec.Command("docker-compose", "ps")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run()
}

func checkServiceConnectivity() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	_, _ = blue.Println("Verificando conectividade dos serviços...")

	// Service Bus Emulator
	_, _ = blue.Println("  - Service Bus Emulator (porta 5672)...")
	if checkPort(5672) {
		_, _ = green.Println("    ✅ Service Bus Emulator está funcionando")
	} else {
		_, _ = yellow.Println("    ⚠️  Service Bus Emulator não está respondendo")
	}

	// Azurite (Azure Storage)
	_, _ = blue.Println("  - Azurite Storage (porta 10000)...")
	if checkHTTPEndpointStart("http://localhost:10000") {
		_, _ = green.Println("    ✅ Azurite Storage está funcionando")
	} else {
		_, _ = yellow.Println("    ⚠️  Azurite Storage não está respondendo")
	}

	// Orion API Mock
	_, _ = blue.Println("  - Orion API (porta 3333)...")
	if checkHTTPEndpointStart("http://localhost:3333") {
		_, _ = green.Println("    ✅ Orion API está funcionando")
	} else {
		_, _ = yellow.Println("    ⚠️  Orion API não está respondendo")
	}

	// Orion Functions
	_, _ = blue.Println("  - Orion Functions (porta 7071)...")
	time.Sleep(10 * time.Second)
	if checkHTTPEndpointStart("http://localhost:7071") {
		_, _ = green.Println("    ✅ Orion Functions está funcionando")
	} else {
		_, _ = yellow.Println("    ⚠️  Orion Functions pode estar ainda inicializando...")
	}
}

func checkPort(port int) bool {
	return utils.CheckPort(port)
}

func checkHTTPEndpointStart(url string) bool {
	return utils.CheckHTTPEndpoint(url)
}

func showFinalInfoStart() {
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	fmt.Println()
	_, _ = green.Println("🎉 Ambiente iniciado com sucesso!")
	fmt.Println()
	_, _ = blue.Println("📋 Serviços disponíveis:")
	fmt.Println("  - Orion Functions: http://localhost:7071")
	fmt.Println("  - Orion API: http://localhost:3333")
	fmt.Println("  - Azurite Storage: http://localhost:10000")
	fmt.Println()
	_, _ = blue.Println("🧪 Para testar as functions:")
	fmt.Println("  - QR Code COB: curl http://localhost:7071/cob/test-id")
	fmt.Println("  - QR Code COBV: curl http://localhost:7071/cobv/test-id")
	fmt.Println("  - Orion API: curl -H 'X-API-Key: FAKE-API-KEY' http://localhost:3333/health")
	fmt.Println()
	_, _ = blue.Println("📝 Logs dos containers:")
	fmt.Println("  - docker-compose logs -f orion-functions")
	fmt.Println("  - docker-compose logs -f emulator")
	fmt.Println()
	_, _ = blue.Println("🛠️  Comandos úteis:")
	fmt.Println("  - orion-dev status        - Ver status dos containers")
	fmt.Println("  - orion-dev test          - Executar testes")
	fmt.Println("  - orion-dev logs          - Ver logs")
	fmt.Println("  - orion-dev stop          - Parar ambiente")
}
