package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// ValidateJSONFile valida se um arquivo contém JSON válido
func ValidateJSONFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	return ValidateJSON(data)
}

// ValidateJSON valida se os dados contêm JSON válido
func ValidateJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("json inválido: %w", err)
	}
	return nil
}

// FormatJSON formata JSON com indentação
func FormatJSON(data []byte) ([]byte, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("json inválido: %w", err)
	}

	return json.MarshalIndent(v, "", "  ")
}

// FormatJSONFile formata um arquivo JSON
func FormatJSONFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	formatted, err := FormatJSON(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, formatted, 0644)
}

// PrettyPrintJSON imprime JSON formatado no terminal
func PrettyPrintJSON(data []byte) error {
	formatted, err := FormatJSON(data)
	if err != nil {
		return err
	}

	fmt.Println(string(formatted))
	return nil
}

// PrettyPrintJSONFile imprime um arquivo JSON formatado no terminal
func PrettyPrintJSONFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	return PrettyPrintJSON(data)
}
