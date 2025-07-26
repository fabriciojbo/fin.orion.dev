package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configurar ambiente inicial",
	Long:  `Configura o ambiente de desenvolvimento Orion com todas as dependências necessárias.`,
	RunE:  runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)

	_, _ = blue.Println("🔧 Configurando ambiente inicial...")
	fmt.Println()

	// Verificar dependências
	if err := checkDependencies(); err != nil {
		return err
	}

	// Configurar arquivos de ambiente
	if err := setupConfigFiles(); err != nil {
		return err
	}

	// Instalar dependências do projeto
	if err := installProjectDeps(); err != nil {
		return err
	}

	// Verificar estrutura do projeto
	if err := checkProjectStructure(); err != nil {
		return err
	}

	// Mostrar informações finais
	showEnvironmentInfo()

	return nil
}

func checkDependencies() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Println("Verificando dependências...")

	deps := []struct {
		name string
		cmd  string
	}{
		{"Docker", "docker"},
		{"Docker Compose", "docker-compose"},
		{"Node.js", "node"},
	}

	for _, dep := range deps {
		if _, err := exec.LookPath(dep.cmd); err != nil {
			_, _ = red.Printf("❌ %s não está instalado\n", dep.name)
			return fmt.Errorf("%s não está instalado", dep.name)
		}
		_, _ = green.Printf("✅ %s encontrado\n", dep.name)
	}

	// Verificar se o Docker está rodando
	if err := exec.Command("docker", "info").Run(); err != nil {
		_, _ = red.Println("❌ Docker não está rodando")
		return fmt.Errorf("docker não está rodando")
	}

	_, _ = green.Println("Todas as dependências verificadas!")
	return nil
}

func setupConfigFiles() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Configurando arquivos de ambiente...")

	// Criar pasta docker/database/certs se não existir
	certsDir := "docker/database/certs"
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		return err
	}
	_, _ = green.Println("Pasta docker/database/certs criada/verificada")

	// Criar pasta docker/service-bus/certs se não existir
	certsDir = "docker/service-bus/certs"
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		return err
	}
	_, _ = green.Println("Pasta docker/service-bus/certs criada/verificada")

	// Criar pasta messages se não existir
	messagesDir := "messages"
	if err := os.MkdirAll(messagesDir, 0755); err != nil {
		return err
	}
	_, _ = green.Println("Pasta messages criada/verificada")

	// Verificar e gerar arquivo .env
	if err := checkAndGenerateEnvFile(); err != nil {
		return err
	}

	// Verificar e gerar arquivo local.settings.json
	if err := checkAndGenerateLocalSettings(); err != nil {
		return err
	}

	// Verificar e copiar Dockerfiles
	if err := copyDockerfiles(); err != nil {
		return err
	}

	_, _ = green.Println("Arquivos de ambiente configurados!")
	return nil
}

func copyDockerfiles() error {
	// Copiar Dockerfile.api para Fin.Orion.API
	apiDockerfile := "../Fin.Orion.API/source/Dockerfile"
	if _, err := os.Stat(apiDockerfile); os.IsNotExist(err) {
		sourceFile := "./docker/container/Dockerfile.api"
		if _, err := os.Stat(sourceFile); err == nil {
			if err := copyFile(sourceFile, apiDockerfile); err != nil {
				return err
			}
		}
	}

	// Copiar Dockerfile.functions para Fin.Orion.Functions
	funcDockerfile := "../Fin.Orion.Functions/source/Dockerfile"
	if _, err := os.Stat(funcDockerfile); os.IsNotExist(err) {
		sourceFile := "./docker/container/Dockerfile.functions"
		if _, err := os.Stat(sourceFile); err == nil {
			if err := copyFile(sourceFile, funcDockerfile); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Criar diretório de destino se não existir
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

func installProjectDeps() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Verificando dependências do projeto...")

	// Verificar se o binário Go foi compilado
	if _, err := os.Stat("bin/orion-dev"); err != nil {
		_, _ = blue.Println("Compilando aplicação Go...")
		cmd := exec.Command("go", "build", "-o", "bin/orion-dev", "cmd/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao compilar aplicação Go: %w", err)
		}
		_, _ = green.Println("Aplicação Go compilada!")
	} else {
		_, _ = green.Println("Aplicação Go já compilada")
	}

	// Verificar certificados
	certFile := "docker/database/certs/server.crt"
	if _, err := os.Stat(certFile); err != nil || os.IsNotExist(err) {
		_, _ = green.Println("Certificados não encontrados. Gerando...")
		if err := generateCertificates(); err != nil {
			return err
		}
		_, _ = green.Println("Certificados gerados com sucesso")
	} else {
		_, _ = green.Println("Certificados encontrados")
	}

	return nil
}

func generateCertificates() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("Gerando certificados PostgreSQL...")

	// Criar diretório se não existir
	certsDir := "docker/database/certs"
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de certificados: %w", err)
	}

	// Gerar chave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("erro ao gerar chave privada: %w", err)
	}

	// Criar template do certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"BR"},
			Province:     []string{"SP"},
			Locality:     []string{"Local"},
			Organization: []string{"Dev"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Válido por 1 ano
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "postgres", "orion-database"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	// Criar certificado
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("erro ao criar certificado: %w", err)
	}

	// Salvar certificado público
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	certPath := filepath.Join(certsDir, "server.crt")
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("erro ao salvar certificado: %w", err)
	}

	// Salvar chave privada
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	keyPath := filepath.Join(certsDir, "server.key")
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("erro ao salvar chave privada: %w", err)
	}

	_, _ = green.Printf("📜 Certificado PostgreSQL salvo em: %s\n", certPath)
	_, _ = green.Printf("🔑 Chave privada PostgreSQL salva em: %s\n", keyPath)

	return nil
}

func checkProjectStructure() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Println("Verificando estrutura do projeto...")

	requiredFiles := []string{
		".env",
		"docker-compose.yml",
		"local.settings.json",
		"docker/database/certs/server.crt",
		"docker/database/certs/server.key",
		"docker/database/init-postgres.sql",
		"docker/database/postgres.conf",
		"docker/service-bus/config.json",
		"docker/container/Dockerfile.api",
		"docker/container/Dockerfile.functions",
		"docker/container/.dockerignore",
	}

	var missingFiles []string
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		_, _ = red.Println("Arquivos ausentes:")
		for _, file := range missingFiles {
			_, _ = red.Printf("  - %s\n", file)
		}
		return fmt.Errorf("arquivos obrigatórios não encontrados")
	}

	_, _ = green.Println("Estrutura do projeto verificada!")
	return nil
}

func showEnvironmentInfo() {
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	fmt.Println()
	_, _ = green.Println("🎉 Ambiente configurado com sucesso!")
	fmt.Println()
	_, _ = blue.Println("📋 Informações do ambiente:")
	fmt.Println("  - Docker: Disponível")
	fmt.Println("  - Docker Compose: Disponível")
	fmt.Println("  - Node.js: Disponível")
	fmt.Println()
	_, _ = blue.Println("🚀 Próximos passos:")
	fmt.Println("  1. Iniciar ambiente: orion-dev start")
	fmt.Println("  2. Verificar status: orion-dev status")
	fmt.Println("  3. Testar functions: orion-dev test")
	fmt.Println()
	_, _ = blue.Println("📚 Comandos úteis:")
	fmt.Println("  - orion-dev help          - Ver todos os comandos")
	fmt.Println("  - orion-dev start         - Iniciar ambiente")
	fmt.Println("  - orion-dev stop          - Parar ambiente")
	fmt.Println("  - orion-dev status        - Ver status")
	fmt.Println("  - orion-dev push-message  - Enviar mensagens")
}

func checkAndGenerateEnvFile() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		_, _ = blue.Println("Arquivo .env não encontrado. Gerando...")

		envContent := `# Configurações básicas
PORT=3333
ENV=HMG
API_KEY=FAKE-API-KEY

# JWT
JWT_SECRET=FAKE-JWT-SECRET
JWT_EXPIRE_MS=3600

# Application Insights
APPLICATIONINSIGHTS_CONNECTION_STRING="InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://brazilsouth-0.in.applicationinsights.azure.com/;LiveEndpoint=https://brazilsouth.livediagnostics.monitor.azure.com/;ApplicationId=00000000-0000-0000-0000-000000000000"

# PostgreSQL
PG_HOST="orion-database"
PG_PORT=5432
PG_USERNAME="postgres"
PG_PASSWORD="postgres"
PG_DATABASE="orion_database"
PG_SCHEMA="orionlocal"

# Pismo
PISMO_URL="https://sandbox.pismolabs.io"
PISMO_SERVER_KEY="YOUR_PISMO_SERVER_KEY"
PISMO_SERVER_SECRET="YOUR_PISMO_SERVER_SECRET"
PISMO_PROGRAM_ID="YOUR_PISMO_PROGRAM_ID"

# PIX
PIX_QRCODE_URL="pix-qrcode-h.magfinancas.com.br"

# Key Vault
KV_RESOURCE_NAME="YOUR_KV_RESOURCE_NAME"
KV_PISMOCERT_NAME="YOUR_KV_PISMOCERT_NAME"

# Contas
PAYMENT_ACCOUNT_SEQ=000000000000
BANK_BRANCH=0000

# MAG IP
MAG_IP_PAYMENT_ACCOUNT_ID="00000000-0000-0000-0000-000000000000"
MAG_IP_PISMO_ACCOUNT_ID=000000000000
MAG_IP_ACCOUNT_BRANCH=0000
MAG_IP_ACCOUNT_NUMBER=00000000

# Cron
DISABLE_CRON=false
CONSULT_BILLET_CRON="0 */1 * * *"

# Celcoin
CELCOIN_URL="YOUR_CELCOIN_URL"
CELCOIN_CLIENT_ID="YOUR_CELCOIN_CLIENT_ID"
CELCOIN_CLIENT_SECRET="YOUR_CELCOIN_CLIENT_SECRET"

# Service Bus
SB_CNT_STR="Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=FAKE-SAS-KEY-VALUE;UseDevelopmentEmulator=true"
SB_QUEUE_NAME="sbq.pismo.onboarding.succeeded"

# Azure Storage
AZURE_STORAGE_BLOB_STRING_CONNECTION="DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://azure-storage:10000/devstoreaccount1;QueueEndpoint=http://azure-storage:10001/devstoreaccount1;TableEndpoint=http://azure-storage:10002/devstoreaccount1;"
AZURE_STORAGE_BLOB_ACCOUNT_UPLOADS_CONTAINER="account-uploads"

# RSA Keys
RSA_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\nYOUR_RSA_PUBLIC_KEY\n-----END PUBLIC KEY-----"
RSA_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYOUR_RSA_PRIVATE_KEY\n-----END PRIVATE KEY-----"

# Magfinancas
MAGFINANCAS_SENSEDIA_BASE_URL="YOUR_MAGFINANCAS_SENSEDIA_BASE_URL"
MAGFINANCAS_SENSEDIA_BEARER_TOKEN="YOUR_MAGFINANCAS_SENSEDIA_BEARER_TOKEN"

# SQL Server Edge & Service Bus Emulator
ACCEPT_EULA="Y"
MSSQL_SA_PASSWORD="YOUR_MSSQL_SA_PASSWORD"
`

		if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo .env: %w", err)
		}
		_, _ = green.Printf("📄 Arquivo .env criado: %s\n", envFile)
	} else {
		_, _ = green.Println("Arquivo .env encontrado")
	}

	return nil
}

func checkAndGenerateLocalSettings() error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	localSettingsFile := "local.settings.json"
	if _, err := os.Stat(localSettingsFile); os.IsNotExist(err) {
		_, _ = blue.Println("Arquivo local.settings.json não encontrado. Gerando...")

		localSettingsContent := `{
  "IsEncrypted": false,
  "Values": {
    "DEBUG": 1,
    "AzureWebJobsFeatureFlags": "EnableWorkerIndexing",
    "APPLICATIONINSIGHTS_CONNECTION_STRING": "InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://brazilsouth-0.in.applicationinsights.azure.com/;LiveEndpoint=https://brazilsouth.livediagnostics.monitor.azure.com/",
    "SB_CONN_STR": "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true",
    "SB_INPUT_QUEUE": "sbq.pismo.transaction.creation",
    "SB_OUTPUT_TOPIC": "sbt.orion.core",
    "PG_HOST": "orion-database",
    "PG_PORT": "5432",
    "PG_USERNAME": "postgres",
    "PG_PASSWORD": "postgres",
    "PG_DATABASE": "orion_database",
    "PG_SCHEMA": "orionlocal,orionlocal,orionlocal,orionlocal",
    "PG_SCHEMA_PIX": "orionlocal",
    "PISMO_URL": "https://sandbox.pismolabs.io",
    "PISMO_SERVER_KEY": "YOUR_PISMO_SERVER_KEY",
    "PISMO_SERVER_SECRET": "YOUR_PISMO_SERVER_SECRET",
    "PISMO_PROGRAM_ID": "YOUR_PISMO_PROGRAM_ID",
    "ORION_URL": "http://orion-api:3333",
    "ORION_API_KEY": "FAKE-API-KEY",
    "PIX_STG_CONN_STR": "DefaultEndpointsProtocol=https;AccountName=YOUR_ACCOUNT_NAME;AccountKey=YOUR_ACCOUNT_KEY;EndpointSuffix=core.windows.net"
  }
}`

		if err := os.WriteFile(localSettingsFile, []byte(localSettingsContent), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo local.settings.json: %w", err)
		}
		_, _ = green.Printf("📄 Arquivo local.settings.json criado: %s\n", localSettingsFile)
	} else {
		_, _ = green.Println("Arquivo local.settings.json encontrado")
	}

	return nil
}
