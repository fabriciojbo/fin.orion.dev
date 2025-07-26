package servicebus

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
)

// Client representa o cliente do Service Bus
type Client struct {
	client *azservicebus.Client
}

// Message representa uma mensagem do Service Bus
type Message struct {
	Body            interface{}            `json:"body"`
	MessageID       string                 `json:"messageId,omitempty"`
	CorrelationID   string                 `json:"correlationId,omitempty"`
	ContentType     string                 `json:"contentType,omitempty"`
	Properties      map[string]interface{} `json:"properties,omitempty"`
	EnqueuedTimeUtc *time.Time             `json:"enqueuedTimeUtc,omitempty"`
	DeliveryCount   int32                  `json:"deliveryCount,omitempty"`
}

// NewClient cria um novo cliente do Service Bus
func NewClient() (*Client, error) {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao carregar .env: %v\n", err)
	}

	connectionString := os.Getenv("SB_CNT_STR")
	if connectionString == "" {
		// Fallback para connection string padrão do emulador
		connectionString = "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=FAKE-SAS-KEY-VALUE"
	}

	fmt.Printf("🔍 Conectando ao Service Bus...\n")

	// Configurar opções do cliente para usar AMQP sem TLS
	clientOptions := &azservicebus.ClientOptions{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		RetryOptions: azservicebus.RetryOptions{
			MaxRetries: 3,
		},
	}

	client, err := azservicebus.NewClientFromConnectionString(connectionString, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente: %w", err)
	}

	fmt.Printf("✅ Cliente criado com sucesso\n")
	return &Client{client: client}, nil
}

// Close fecha a conexão do cliente
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close(context.TODO())
	}
	return nil
}

// TestConnection testa a conectividade básica com o Service Bus
func (c *Client) TestConnection() error {
	// Tentar criar um sender temporário para testar a conexão
	sender, err := c.client.NewSender("test-queue", nil)
	if err != nil {
		return fmt.Errorf("erro ao testar conexão: %w", err)
	}
	defer func() { _ = sender.Close(context.TODO()) }()

	fmt.Printf("✅ Conexão com Service Bus estabelecida com sucesso\n")
	return nil
}

// SendMessageToQueue envia uma mensagem para uma fila
func (c *Client) SendMessageToQueue(queueName string, message *Message) error {
	fmt.Printf("🔧 Criando sender para fila: %s\n", queueName)

	// Converter mensagem para o formato do Azure Service Bus
	body, err := json.Marshal(message.Body)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}
	fmt.Printf("✅ Mensagem serializada (%d bytes)\n", len(body))

	sender, err := c.client.NewSender(queueName, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar sender: %w", err)
	}
	defer func() { _ = sender.Close(context.TODO()) }()
	fmt.Printf("✅ Sender criado com sucesso\n")

	sbMessage := &azservicebus.Message{
		Body: body,
	}

	// Adicionar MessageID se não estiver vazio
	if message.MessageID != "" {
		sbMessage.MessageID = &message.MessageID
		fmt.Printf("📝 MessageID definido: %s\n", message.MessageID)
	}

	// Adicionar CorrelationID se não estiver vazio
	if message.CorrelationID != "" {
		sbMessage.CorrelationID = &message.CorrelationID
		fmt.Printf("🔗 CorrelationID definido: %s\n", message.CorrelationID)
	}

	// Adicionar ContentType se não estiver vazio
	if message.ContentType != "" {
		sbMessage.ContentType = &message.ContentType
		fmt.Printf("📋 ContentType definido: %s\n", message.ContentType)
	}

	// Adicionar propriedades se existirem
	if message.Properties != nil {
		sbMessage.ApplicationProperties = message.Properties
		fmt.Printf("🏷️  Propriedades adicionadas: %d\n", len(message.Properties))
	}

	fmt.Printf("📤 Enviando mensagem para Service Bus...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = sender.SendMessage(ctx, sbMessage, nil)
	cancel()

	if err != nil {
		return fmt.Errorf("erro ao enviar mensagem: %w", err)
	}

	return nil
}

// ReceiveMessagesFromQueue recebe mensagens de uma fila
func (c *Client) ReceiveMessagesFromQueue(queueName string, maxMessages int) ([]*Message, error) {
	receiver, err := c.client.NewReceiverForQueue(queueName, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar receiver: %w", err)
	}
	defer func() { _ = receiver.Close(context.TODO()) }()

	messages, err := receiver.ReceiveMessages(context.TODO(), maxMessages, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao receber mensagens: %w", err)
	}

	var result []*Message
	for _, msg := range messages {
		message := &Message{
			MessageID:       msg.MessageID,
			EnqueuedTimeUtc: msg.EnqueuedTime,
			DeliveryCount:   int32(msg.DeliveryCount),
		}

		// Adicionar CorrelationID se não for nil
		if msg.CorrelationID != nil {
			message.CorrelationID = *msg.CorrelationID
		}

		// Adicionar ContentType se não for nil
		if msg.ContentType != nil {
			message.ContentType = *msg.ContentType
		}

		// Deserializar body
		if err := json.Unmarshal(msg.Body, &message.Body); err != nil {
			message.Body = string(msg.Body)
		}

		// Adicionar propriedades
		if msg.ApplicationProperties != nil {
			message.Properties = make(map[string]interface{})
			for k, v := range msg.ApplicationProperties {
				message.Properties[k] = v
			}
		}

		result = append(result, message)

		// Complete a mensagem
		if err := receiver.CompleteMessage(context.TODO(), msg, nil); err != nil {
			return nil, fmt.Errorf("erro ao completar mensagem: %w", err)
		}
	}

	return result, nil
}

// SendMessageToTopic envia uma mensagem para um tópico
func (c *Client) SendMessageToTopic(topicName string, message *Message) error {
	sender, err := c.client.NewSender(topicName, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar sender: %w", err)
	}
	defer func() { _ = sender.Close(context.TODO()) }()

	// Converter mensagem para o formato do Azure Service Bus
	body, err := json.Marshal(message.Body)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	sbMessage := &azservicebus.Message{
		Body: body,
	}

	// Adicionar MessageID se não estiver vazio
	if message.MessageID != "" {
		sbMessage.MessageID = &message.MessageID
	}

	// Adicionar CorrelationID se não estiver vazio
	if message.CorrelationID != "" {
		sbMessage.CorrelationID = &message.CorrelationID
	}

	// Adicionar ContentType se não estiver vazio
	if message.ContentType != "" {
		sbMessage.ContentType = &message.ContentType
	}

	// Adicionar propriedades se existirem
	if message.Properties != nil {
		sbMessage.ApplicationProperties = message.Properties
	}

	return sender.SendMessage(context.TODO(), sbMessage, nil)
}

// ReceiveMessagesFromTopic recebe mensagens de uma subscription de tópico
func (c *Client) ReceiveMessagesFromTopic(topicName, subscriptionName string, maxMessages int) ([]*Message, error) {
	receiver, err := c.client.NewReceiverForSubscription(topicName, subscriptionName, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar receiver: %w", err)
	}
	defer func() { _ = receiver.Close(context.TODO()) }()

	messages, err := receiver.ReceiveMessages(context.TODO(), maxMessages, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao receber mensagens: %w", err)
	}

	var result []*Message
	for _, msg := range messages {
		message := &Message{
			MessageID:       msg.MessageID,
			EnqueuedTimeUtc: msg.EnqueuedTime,
			DeliveryCount:   int32(msg.DeliveryCount),
		}

		// Adicionar CorrelationID se não for nil
		if msg.CorrelationID != nil {
			message.CorrelationID = *msg.CorrelationID
		}

		// Adicionar ContentType se não for nil
		if msg.ContentType != nil {
			message.ContentType = *msg.ContentType
		}

		// Deserializar body
		if err := json.Unmarshal(msg.Body, &message.Body); err != nil {
			message.Body = string(msg.Body)
		}

		// Adicionar propriedades
		if msg.ApplicationProperties != nil {
			message.Properties = make(map[string]interface{})
			for k, v := range msg.ApplicationProperties {
				message.Properties[k] = v
			}
		}

		result = append(result, message)

		// Complete a mensagem
		if err := receiver.CompleteMessage(context.TODO(), msg, nil); err != nil {
			return nil, fmt.Errorf("erro ao completar mensagem: %w", err)
		}
	}

	return result, nil
}
