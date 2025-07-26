package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fin.orion.dev/internal/proxy"
	"fin.orion.dev/internal/servicebus"
	"fin.orion.dev/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Lista de filas válidas
var validQueues = []string{
	"sbq.pismo.onboarding.succeeded",
	"sbq.pismo.transaction.creation",
	"sbq.pismo.pix.transaction.in",
	"sbq.pismo.all",
	"sbq.orion.pixqrcode.persist",
	"sbq.orion.transaction.chained",
	"sbq.orion.billet-payment.verify",
	"sbq.pismo.authorization.cancelation",
	"sbq.pismo.ted.transaction",
	"sbq.pix.recurrence.payment.order.failure",
}

// Comando para verificar mensagens
var checkMessagesCmd = &cobra.Command{
	Use:   "check-messages",
	Short: "Verificar mensagens do Service Bus",
	Long:  `Verifica mensagens das filas e tópicos do Service Bus.`,
	RunE:  runCheckMessages,
}

// Comando para enviar mensagens
var pushMessageCmd = &cobra.Command{
	Use:   "push-message [queue] [file]",
	Short: "Enviar mensagem para fila",
	Long:  `Envia uma mensagem JSON para uma fila específica.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runPushMessage,
}

// Comando para verificar filas
var checkQueueCmd = &cobra.Command{
	Use:   "check-queue [queue]",
	Short: "Verificar mensagens da fila",
	Long:  `Verifica mensagens de uma fila específica.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCheckQueue,
}

// Comando para verificar tópicos
var checkTopicCmd = &cobra.Command{
	Use:   "check-topic [subscription]",
	Short: "Verificar mensagens do tópico",
	Long:  `Verifica mensagens de um tópico específico.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCheckTopic,
}

// Comando para listar
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar recursos",
	Long:  `Lista filas, arquivos JSON e outros recursos disponíveis.`,
	RunE:  runList,
}

// Comando para teste rápido de mensagem
var testMessageCmd = &cobra.Command{
	Use:   "test-message",
	Short: "Teste rápido de envio de mensagem",
	Long:  `Executa um teste rápido de envio de mensagem para verificar conectividade.`,
	RunE:  runTestMessage,
}

// Comando para enviar mensagem de teste para fila
var sendQueueCmd = &cobra.Command{
	Use:   "send-queue [queue] [type]",
	Short: "Enviar mensagem de teste para fila",
	Long:  `Envia uma mensagem de teste para uma fila específica.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runSendQueue,
}

// Comando para enviar mensagem JSON para fila
var sendJsonCmd = &cobra.Command{
	Use:   "send-json [queue] [file]",
	Short: "Enviar mensagem JSON para fila",
	Long:  `Envia uma mensagem JSON para uma fila específica.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runSendJson,
}

// Comando para listar filas
var listQueuesCmd = &cobra.Command{
	Use:   "list-queues",
	Short: "Listar filas disponíveis",
	Long:  `Lista todas as filas disponíveis no Service Bus.`,
	RunE:  runListQueues,
}

// Comando para listar mensagens
var listMessagesCmd = &cobra.Command{
	Use:   "list-messages",
	Short: "Listar arquivos JSON disponíveis",
	Long:  `Lista todos os arquivos JSON disponíveis na pasta messages.`,
	RunE:  runListMessages,
}

// Comando para validar JSON
var validateJsonCmd = &cobra.Command{
	Use:   "validate-json [file]",
	Short: "Validar arquivo JSON",
	Long:  `Valida se um arquivo contém JSON válido.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runValidateJson,
}

// Comando para formatar JSON
var formatJsonCmd = &cobra.Command{
	Use:   "format-json [file]",
	Short: "Formatar arquivo JSON",
	Long:  `Formata um arquivo JSON com indentação adequada.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runFormatJson,
}

// Comando para mostrar JSON formatado
var showJsonCmd = &cobra.Command{
	Use:   "show-json [file]",
	Short: "Mostrar JSON formatado",
	Long:  `Mostra o conteúdo de um arquivo JSON formatado no terminal.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runShowJson,
}

// Comando para iniciar proxy do Service Bus
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Iniciar proxy do Service Bus",
	Long:  `Inicia um proxy TCP que redireciona conexões da porta 5671 para 5672.`,
	RunE:  runProxy,
}

func runCheckMessages(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	_, _ = blue.Println("📨 Verificando mensagens do Service Bus...")

	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()

	// Verificar status do ambiente
	checkEnvironmentStatus()

	return nil
}

func runPushMessage(cmd *cobra.Command, args []string) error {
	queueName := args[0]
	jsonFile := args[1]

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("📤 Enviando mensagem para fila: %s\n", queueName)
	_, _ = blue.Printf("📁 Arquivo: %s\n", jsonFile)
	fmt.Println()

	// Validar nome da fila
	if !isValidQueue(queueName) {
		_, _ = red.Printf("❌ Fila '%s' não é válida\n", queueName)
		_, _ = blue.Println("Filas válidas:")
		for _, queue := range validQueues {
			_, _ = blue.Printf("  - %s\n", queue)
		}
		return fmt.Errorf("fila inválida")
	}

	// Carregar mensagem do arquivo JSON
	_, _ = blue.Println("📄 Carregando mensagem do arquivo...")
	message, err := loadMessageFromFile(jsonFile)
	if err != nil {
		return fmt.Errorf("erro ao carregar arquivo: %w", err)
	}
	_, _ = green.Println("✅ Mensagem carregada com sucesso")

	// Conectar ao Service Bus
	_, _ = blue.Println("🔗 Conectando ao Service Bus...")
	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()
	_, _ = green.Println("✅ Conectado ao Service Bus")

	// Enviar mensagem
	_, _ = blue.Println("📤 Enviando mensagem...")
	if err := client.SendMessageToQueue(queueName, message); err != nil {
		return fmt.Errorf("erro ao enviar mensagem: %w", err)
	}

	_, _ = green.Println("✅ Mensagem enviada com sucesso!")
	return nil
}

func runCheckQueue(cmd *cobra.Command, args []string) error {
	queueName := args[0]

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("🔍 Verificando mensagens da fila '%s'...\n", queueName)
	fmt.Println()

	// Validar nome da fila
	if !isValidQueue(queueName) {
		_, _ = red.Printf("❌ Fila '%s' não é válida\n", queueName)
		_, _ = blue.Println("Filas válidas:")
		for _, queue := range validQueues {
			_, _ = blue.Printf("  - %s\n", queue)
		}
		return fmt.Errorf("fila inválida")
	}

	// Conectar ao Service Bus
	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()

	// Receber mensagens
	messages, err := client.ReceiveMessagesFromQueue(queueName, 10)
	if err != nil {
		return fmt.Errorf("erro ao receber mensagens: %w", err)
	}

	if len(messages) == 0 {
		_, _ = blue.Println("ℹ️  Nenhuma mensagem encontrada na fila")
	} else {
		_, _ = green.Printf("📨 Encontradas %d mensagem(ns):\n", len(messages))
		fmt.Println()

		for i, message := range messages {
			_, _ = blue.Printf("--- Mensagem %d ---\n", i+1)
			fmt.Printf("ID: %s\n", message.MessageID)
			fmt.Printf("Correlation ID: %s\n", message.CorrelationID)
			fmt.Printf("Content Type: %s\n", message.ContentType)
			if message.EnqueuedTimeUtc != nil {
				fmt.Printf("Timestamp: %s\n", message.EnqueuedTimeUtc.Format(time.RFC3339))
			}
			fmt.Printf("Delivery Count: %d\n", message.DeliveryCount)
			fmt.Println("Body:")

			bodyJSON, _ := json.MarshalIndent(message.Body, "", "  ")
			fmt.Println(string(bodyJSON))
			fmt.Println()
		}
	}

	_, _ = green.Println("✅ Verificação concluída")
	return nil
}

func runCheckTopic(cmd *cobra.Command, args []string) error {
	subscription := "subscription.orion.core"
	if len(args) > 0 {
		subscription = args[0]
	}

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🔍 Verificando mensagens do tópico 'sbt.orion.core'...")
	_, _ = blue.Printf("📋 Subscription: %s\n", subscription)
	fmt.Println()

	// Conectar ao Service Bus
	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()

	// Receber mensagens
	messages, err := client.ReceiveMessagesFromTopic("sbt.orion.core", subscription, 10)
	if err != nil {
		return fmt.Errorf("erro ao receber mensagens: %w", err)
	}

	if len(messages) == 0 {
		_, _ = blue.Println("ℹ️  Nenhuma mensagem encontrada no tópico")
	} else {
		_, _ = green.Printf("📨 Encontradas %d mensagem(ns):\n", len(messages))
		fmt.Println()

		for i, message := range messages {
			_, _ = blue.Printf("--- Mensagem %d ---\n", i+1)
			fmt.Printf("ID: %s\n", message.MessageID)
			fmt.Printf("Correlation ID: %s\n", message.CorrelationID)
			fmt.Printf("Content Type: %s\n", message.ContentType)
			if message.EnqueuedTimeUtc != nil {
				fmt.Printf("Timestamp: %s\n", message.EnqueuedTimeUtc.Format(time.RFC3339))
			}
			fmt.Println("Body:")

			bodyJSON, _ := json.MarshalIndent(message.Body, "", "  ")
			fmt.Println(string(bodyJSON))
			fmt.Println()
		}
	}

	_, _ = green.Println("✅ Verificação concluída")
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("📋 Listando recursos disponíveis...")
	fmt.Println()

	// Listar filas
	_, _ = blue.Println("📨 Filas disponíveis:")
	for i, queue := range validQueues {
		_, _ = green.Printf("  %d. %s\n", i+1, queue)
	}
	fmt.Println()

	// Listar arquivos JSON
	_, _ = blue.Println("📁 Arquivos JSON disponíveis:")
	files, err := listJSONFiles()
	if err != nil {
		return fmt.Errorf("erro ao listar arquivos: %w", err)
	}

	if len(files) == 0 {
		_, _ = blue.Println("  Nenhum arquivo JSON encontrado na pasta 'messages'")
	} else {
		for _, file := range files {
			_, _ = green.Printf("  - %s\n", file)
		}
	}

	return nil
}

func runTestMessage(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🔗 Testando conexão com Service Bus...")
	fmt.Println()

	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()

	queueName := "test-queue" // Nome da fila de teste
	message := &servicebus.Message{
		Body:          map[string]string{"message": "Hello from test-message!"},
		MessageID:     fmt.Sprintf("test-%d", time.Now().Unix()),
		CorrelationID: fmt.Sprintf("corr-%d", time.Now().Unix()),
		ContentType:   "application/json",
	}

	if err := client.SendMessageToQueue(queueName, message); err != nil {
		return fmt.Errorf("erro ao enviar mensagem de teste: %w", err)
	}

	_, _ = green.Printf("✅ Mensagem de teste enviada para a fila '%s'\n", queueName)
	return nil
}

func runSendQueue(cmd *cobra.Command, args []string) error {
	queueName := args[0]
	messageType := "simple"
	if len(args) > 1 {
		messageType = args[1]
	}

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("📤 Enviando mensagem de teste para fila '%s' (tipo: %s)...\n", queueName, messageType)
	fmt.Println()

	// Validar nome da fila
	if !isValidQueue(queueName) {
		_, _ = red.Printf("❌ Fila '%s' não é válida\n", queueName)
		_, _ = blue.Println("Filas válidas:")
		for _, queue := range validQueues {
			_, _ = blue.Printf("  - %s\n", queue)
		}
		return fmt.Errorf("fila inválida")
	}

	// Conectar ao Service Bus
	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()

	// Criar mensagem baseada no tipo
	var messageBody interface{}
	switch messageType {
	case "simple":
		messageBody = map[string]interface{}{
			"type":      "test",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"message": "Mensagem de teste simples para fila",
				"random":  time.Now().Unix(),
				"queue":   queueName,
			},
		}
	case "pix-recurrence":
		messageBody = map[string]interface{}{
			"type":      "pix.recurrence.payment.order",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"orderId":         fmt.Sprintf("order-%d", time.Now().Unix()),
				"amount":          100.5,
				"description":     "Teste PIX Recurrence para fila",
				"recurrenceType":  "monthly",
				"nextPaymentDate": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
				"queue":           queueName,
			},
		}
	case "transaction":
		messageBody = map[string]interface{}{
			"type":      "transaction.created",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"transactionId": fmt.Sprintf("txn-%d", time.Now().Unix()),
				"amount":        250.0,
				"currency":      "BRL",
				"description":   "Teste de transação para fila",
				"status":        "pending",
				"queue":         queueName,
			},
		}
	case "pismo-transaction":
		messageBody = map[string]interface{}{
			"type":      "pismo.transaction.creation",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"transactionId": fmt.Sprintf("pismo-txn-%d", time.Now().Unix()),
				"amount":        150.0,
				"currency":      "BRL",
				"description":   "Teste transação Pismo",
				"status":        "created",
				"accountId":     "acc-123456",
				"queue":         queueName,
			},
		}
	case "authorization-cancelled":
		messageBody = map[string]interface{}{
			"type":      "pismo.authorization.cancelled",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"authorizationId": fmt.Sprintf("auth-%d", time.Now().Unix()),
				"reason":          "user_cancelled",
				"description":     "Autorização cancelada pelo usuário",
				"queue":           queueName,
			},
		}
	default:
		messageBody = map[string]interface{}{
			"type":      "test",
			"timestamp": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"message": fmt.Sprintf("Mensagem de teste (tipo: %s)", messageType),
				"queue":   queueName,
			},
		}
	}

	// Enviar mensagem
	message := &servicebus.Message{
		Body:          messageBody,
		MessageID:     fmt.Sprintf("send-queue-%s-%d", messageType, time.Now().Unix()),
		CorrelationID: fmt.Sprintf("corr-%d", time.Now().Unix()),
		ContentType:   "application/json",
	}

	if err := client.SendMessageToQueue(queueName, message); err != nil {
		return fmt.Errorf("erro ao enviar mensagem de teste: %w", err)
	}

	_, _ = green.Printf("✅ Mensagem de teste enviada para a fila '%s'\n", queueName)
	return nil
}

func runSendJson(cmd *cobra.Command, args []string) error {
	queueName := args[0]
	jsonFile := args[1]

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("📤 Enviando mensagem JSON para fila '%s'...\n", queueName)
	_, _ = blue.Printf("📁 Arquivo: %s\n", jsonFile)
	fmt.Println()

	// Validar nome da fila
	if !isValidQueue(queueName) {
		_, _ = red.Printf("❌ Fila '%s' não é válida\n", queueName)
		_, _ = blue.Println("Filas válidas:")
		for _, queue := range validQueues {
			_, _ = blue.Printf("  - %s\n", queue)
		}
		return fmt.Errorf("fila inválida")
	}

	// Carregar mensagem do arquivo JSON
	message, err := loadMessageFromFile(jsonFile)
	if err != nil {
		return fmt.Errorf("erro ao carregar arquivo: %w", err)
	}

	// Conectar ao Service Bus
	client, err := servicebus.NewClient()
	if err != nil {
		return fmt.Errorf("erro ao conectar ao service bus: %w", err)
	}
	defer func() { _ = client.Close() }()

	// Enviar mensagem
	if err := client.SendMessageToQueue(queueName, message); err != nil {
		return fmt.Errorf("erro ao enviar mensagem JSON: %w", err)
	}

	_, _ = green.Printf("✅ Mensagem JSON enviada para a fila '%s'\n", queueName)
	return nil
}

func runListQueues(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("📋 Listando filas disponíveis no Service Bus...")
	fmt.Println()

	// Como o client não tem método ListQueues, vamos listar as filas válidas
	_, _ = green.Printf("📨 Filas válidas (%d):\n", len(validQueues))
	fmt.Println()
	for i, queue := range validQueues {
		_, _ = green.Printf("  %d. %s\n", i+1, queue)
	}

	return nil
}

func runListMessages(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("📋 Listando arquivos JSON disponíveis...")
	fmt.Println()

	files, err := listJSONFiles()
	if err != nil {
		return fmt.Errorf("erro ao listar arquivos: %w", err)
	}

	if len(files) == 0 {
		_, _ = blue.Println("ℹ️  Nenhum arquivo JSON encontrado na pasta 'messages'")
	} else {
		_, _ = green.Printf("📁 Encontrados %d arquivo(s):\n", len(files))
		fmt.Println()
		for i, file := range files {
			_, _ = green.Printf("  %d. %s\n", i+1, file)
		}
	}

	return nil
}

func runValidateJson(cmd *cobra.Command, args []string) error {
	jsonFile := args[0]

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("📄 Validando arquivo JSON: %s\n", jsonFile)
	fmt.Println()

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		_, _ = red.Printf("❌ JSON inválido: %s\n", err.Error())
		return fmt.Errorf("json inválido")
	}

	_, _ = green.Println("✅ JSON válido!")
	return nil
}

func runFormatJson(cmd *cobra.Command, args []string) error {
	jsonFile := args[0]

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("📄 Formatando arquivo JSON: %s\n", jsonFile)
	fmt.Println()

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		_, _ = red.Printf("❌ JSON inválido para formatação: %s\n", err.Error())
		return fmt.Errorf("json inválido para formatação")
	}

	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		_, _ = red.Printf("❌ Erro ao formatar JSON: %s\n", err.Error())
		return fmt.Errorf("erro ao formatar json")
	}

	_, _ = green.Println(string(formattedJSON))
	return nil
}

func runShowJson(cmd *cobra.Command, args []string) error {
	jsonFile := args[0]

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Printf("📄 Mostrando conteúdo do arquivo JSON: %s\n", jsonFile)
	fmt.Println()

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		_, _ = red.Printf("❌ JSON inválido para exibição: %s\n", err.Error())
		return fmt.Errorf("json inválido para exibição")
	}

	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		_, _ = red.Printf("❌ Erro ao formatar JSON para exibição: %s\n", err.Error())
		return fmt.Errorf("erro ao formatar json para exibição")
	}

	_, _ = green.Println(string(formattedJSON))
	return nil
}

func runProxy(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("🚀 Iniciando proxy do Service Bus...")
	_, _ = blue.Println("📡 Redirecionando porta 5671 -> 5672")
	fmt.Println()

	// Iniciar o proxy
	if err := proxy.RunProxy(); err != nil {
		return fmt.Errorf("erro ao iniciar proxy: %w", err)
	}

	_, _ = green.Println("✅ Proxy iniciado com sucesso")
	return nil
}

func isValidQueue(queueName string) bool {
	for _, queue := range validQueues {
		if queue == queueName {
			return true
		}
	}
	return false
}

func loadMessageFromFile(filename string) (*servicebus.Message, error) {
	filePath := filepath.Join("messages", filename)

	// Validar JSON antes de carregar
	if err := utils.ValidateJSONFile(filePath); err != nil {
		return nil, fmt.Errorf("erro de validação JSON: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("arquivo não encontrado: %s", filePath)
	}

	var body interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("erro de sintaxe JSON: %w", err)
	}

	message := &servicebus.Message{
		Body:          body,
		MessageID:     fmt.Sprintf("json-%d", time.Now().Unix()),
		CorrelationID: fmt.Sprintf("corr-%d", time.Now().Unix()),
		ContentType:   "application/json",
	}

	return message, nil
}

func listJSONFiles() ([]string, error) {
	files, err := os.ReadDir("messages")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var jsonFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}

	return jsonFiles, nil
}

func checkEnvironmentStatus() {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	_, _ = blue.Println("Verificando status do ambiente...")
	fmt.Println()

	// Verificar Service Bus
	if utils.CheckPort(5672) {
		_, _ = green.Println("✅ Service Bus Emulator (porta 5672)")
	} else {
		_, _ = red.Println("❌ Service Bus Emulator (porta 5672)")
	}

	// Verificar Orion Functions
	if utils.CheckHTTPEndpoint("http://localhost:7071") {
		_, _ = green.Println("✅ Orion Functions (porta 7071)")
	} else {
		_, _ = red.Println("⚠️  Orion Functions (porta 7071)")
	}

	// Verificar Orion API
	if utils.CheckHTTPEndpoint("http://localhost:3333") {
		_, _ = green.Println("✅ Orion API (porta 3333)")
	} else {
		_, _ = red.Println("⚠️  Orion API (porta 3333)")
	}
}
