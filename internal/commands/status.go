package commands

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Ver status dos containers",
	Long:  `Verifica o status dos containers e conectividade dos serviços.`,
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Println("📊 Status dos containers...")
	fmt.Println()

	// Verificar status dos containers
	checkContainerStatusStatus()

	fmt.Println()
	_, _ = blue.Println("🏥 Verificando saúde dos serviços...")

	// Azure Service Bus
	if checkHTTPEndpointStatus("http://localhost:5300/health") {
		_, _ = green.Println("✅ Azure Service Bus: OK")
	} else {
		_, _ = red.Println("❌ Azure Service Bus: ERRO")
	}

	// Azure Storage
	if checkHTTPEndpointStatus("http://localhost:10000") {
		_, _ = green.Println("✅ Azure Storage: OK")
	} else {
		_, _ = red.Println("❌ Azure Storage: ERRO")
	}

	// PostgreSQL
	if checkPostgreSQL() {
		_, _ = green.Println("✅ PostgreSQL: OK")
	} else {
		_, _ = red.Println("❌ PostgreSQL: ERRO")
	}

	// Orion API
	if checkHTTPEndpointStatus("http://localhost:3333") {
		_, _ = green.Println("✅ Orion API: OK")
	} else {
		_, _ = red.Println("❌ Orion API: ERRO")
	}

	// Orion Functions
	if checkHTTPEndpointStatus("http://localhost:7071") {
		_, _ = green.Println("✅ Orion Functions: OK")
	} else {
		_, _ = red.Println("❌ Orion Functions: ERRO")
	}

	return nil
}

func checkContainerStatusStatus() {
	cmd := exec.Command("docker-compose", "ps")
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run()
}

func checkHTTPEndpointStatus(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode < 500
}

func checkPostgreSQL() bool {
	cmd := exec.Command("docker-compose", "exec", "-T", "postgres", "pg_isready", "-U", "postgres")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}
