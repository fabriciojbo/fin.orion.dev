package tests

import (
	"os"
	"strings"
	"testing"

	"fin.orion.dev/internal/utils"
)

func TestGetVersion(t *testing.T) {
	// Criar arquivo .version temporário para teste
	testVersion := "1.0.4"
	err := os.WriteFile(".version", []byte(testVersion), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo .version para teste: %v", err)
	}

	// Limpar arquivo após o teste
	defer os.Remove(".version")

	// Teste com arquivo .version existente
	version, err := utils.GetVersion()
	if err != nil {
		t.Fatalf("Erro ao ler versão: %v", err)
	}

	if version == "" {
		t.Error("Versão não pode estar vazia")
	}

	if version != testVersion {
		t.Errorf("Esperado versão %s, got: %s", testVersion, version)
	}

	// Verificar se a versão tem formato válido (número.número.número)
	if len(version) < 3 {
		t.Errorf("Versão deve ter pelo menos 3 caracteres, got: %s", version)
	}
}

func TestGetVersionOrUnknown(t *testing.T) {
	// Teste normal - criar arquivo .version temporário
	testVersion := "2.0.0"
	err := os.WriteFile(".version", []byte(testVersion), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo .version para teste: %v", err)
	}

	// Limpar arquivo após o teste
	defer os.Remove(".version")

	version := utils.GetVersionOrUnknown()
	if version == "" {
		t.Error("Versão não pode estar vazia")
	}

	if version != testVersion {
		t.Errorf("Esperado versão %s, got: %s", testVersion, version)
	}

	// Teste quando arquivo não existe
	os.Remove(".version")

	versionUnknown := utils.GetVersionOrUnknown()
	if versionUnknown != "unknown" {
		t.Errorf("Esperado 'unknown', got: %s", versionUnknown)
	}
}

func TestGetVersionWithEmptyFile(t *testing.T) {
	// Criar arquivo .version temporário vazio
	err := os.WriteFile(".version", []byte(""), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário vazio: %v", err)
	}

	// Limpar arquivo após o teste
	defer os.Remove(".version")

	// Testar com arquivo vazio
	_, err = utils.GetVersion()
	if err == nil {
		t.Error("Esperado erro para arquivo vazio")
	}

	// Verificar se a mensagem de erro contém "vazio"
	if err != nil && !strings.Contains(err.Error(), "vazio") {
		t.Errorf("Mensagem de erro deve conter 'vazio', got: %v", err)
	}
}
