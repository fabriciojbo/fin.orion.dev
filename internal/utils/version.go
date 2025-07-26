package utils

import (
	"fmt"
	"os"
	"strings"
)

// GetVersion lê a versão do arquivo .version
func GetVersion() (string, error) {
	content, err := os.ReadFile(".version")
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo .version: %w", err)
	}

	version := strings.TrimSpace(string(content))
	if version == "" {
		return "", fmt.Errorf("arquivo .version está vazio")
	}

	return version, nil
}

// GetVersionOrUnknown retorna a versão ou "unknown" se não conseguir ler
func GetVersionOrUnknown() string {
	version, err := GetVersion()
	if err != nil {
		return "unknown"
	}
	return version
}
