package proxy

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// ServiceBusProxy é um proxy TCP que redireciona conexões
type ServiceBusProxy struct {
	listenAddr string
	targetAddr string
	listener   net.Listener
	tlsConfig  *tls.Config
	ctx        context.Context
	cancel     context.CancelFunc
}

// generateCertificate cria um certificado auto-assinado para localhost
func generateCertificate() (*tls.Certificate, error) {
	// Gerar chave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar chave privada: %w", err)
	}

	// Criar template do certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Service Bus Emulator"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Válido por 1 ano
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	// Criar certificado
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar certificado: %w", err)
	}

	// Criar certificado TLS
	cert := &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	return cert, nil
}

// saveCertificate salva o certificado na pasta docker/service-bus/certs
func saveCertificate(cert *tls.Certificate) error {
	certsDir := "docker/service-bus/certs"

	// Criar diretório se não existir
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de certificados: %w", err)
	}

	// Salvar certificado público
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Certificate[0],
	})

	certPath := filepath.Join(certsDir, "servicebus-proxy.crt")
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return fmt.Errorf("erro ao salvar certificado: %w", err)
	}

	// Salvar chave privada
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cert.PrivateKey.(*rsa.PrivateKey)),
	})

	keyPath := filepath.Join(certsDir, "servicebus-proxy.key")
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return fmt.Errorf("erro ao salvar chave privada: %w", err)
	}

	log.Printf("📜 Certificado salvo em: %s", certPath)
	log.Printf("🔑 Chave privada salva em: %s", keyPath)

	return nil
}

// NewServiceBusProxy cria um novo proxy
func NewServiceBusProxy(listenPort, targetPort int) (*ServiceBusProxy, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Gerar certificado TLS
	cert, err := generateCertificate()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("erro ao gerar certificado: %w", err)
	}

	// Salvar certificado
	if err := saveCertificate(cert); err != nil {
		cancel()
		return nil, fmt.Errorf("erro ao salvar certificado: %w", err)
	}

	// Configurar TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
	}

	return &ServiceBusProxy{
		listenAddr: fmt.Sprintf(":%d", listenPort),
		targetAddr: fmt.Sprintf("localhost:%d", targetPort),
		tlsConfig:  tlsConfig,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Start inicia o proxy
func (p *ServiceBusProxy) Start() error {
	var err error
	p.listener, err = tls.Listen("tcp", p.listenAddr, p.tlsConfig)
	if err != nil {
		return fmt.Errorf("erro ao iniciar listener TLS: %w", err)
	}

	log.Printf("🚀 Proxy TLS iniciado: %s -> %s", p.listenAddr, p.targetAddr)

	for {
		select {
		case <-p.ctx.Done():
			return nil
		default:
			conn, err := p.listener.Accept()
			if err != nil {
				log.Printf("❌ Erro ao aceitar conexão: %v", err)
				continue
			}

			go p.handleConnection(conn)
		}
	}
}

// Stop para o proxy
func (p *ServiceBusProxy) Stop() error {
	p.cancel()
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}

// handleConnection gerencia uma conexão individual
func (p *ServiceBusProxy) handleConnection(clientConn net.Conn) {
	defer func() { _ = clientConn.Close() }()

	log.Printf("🔗 Nova conexão de: %s", clientConn.RemoteAddr())

	// Conectar ao target
	targetConn, err := net.DialTimeout("tcp", p.targetAddr, 5*time.Second)
	if err != nil {
		log.Printf("❌ Erro ao conectar ao target %s: %v", p.targetAddr, err)
		return
	}
	defer func() { _ = targetConn.Close() }()

	log.Printf("✅ Conectado ao target: %s", p.targetAddr)

	// Criar canais para sincronização
	errChan := make(chan error, 2)

	// Copiar dados do client para target
	go func() {
		_, err := io.Copy(targetConn, clientConn)
		errChan <- err
	}()

	// Copiar dados do target para client
	go func() {
		_, err := io.Copy(clientConn, targetConn)
		errChan <- err
	}()

	// Aguardar erro ou fechamento
	select {
	case err := <-errChan:
		if err != nil && err != io.EOF {
			log.Printf("⚠️  Erro na transferência de dados: %v", err)
		}
	case <-p.ctx.Done():
		log.Printf("🛑 Proxy parando...")
	}

	log.Printf("🔌 Conexão fechada: %s", clientConn.RemoteAddr())
}

// RunProxy inicia o proxy do Service Bus
func RunProxy() error {
	proxy, err := NewServiceBusProxy(5671, 5672)
	if err != nil {
		return fmt.Errorf("erro ao criar proxy: %w", err)
	}

	// Configurar signal handling para graceful shutdown
	go func() {
		<-proxy.ctx.Done()
		log.Printf("🛑 Recebido sinal de parada")
	}()

	log.Printf("🚀 Iniciando Service Bus Proxy (5671 -> 5672)")
	return proxy.Start()
}
