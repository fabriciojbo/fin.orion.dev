# 🧪 Makefile para Testes do Fin.Orion.Dev
# Autor: Fin.Orion.Dev Team

.PHONY: help test test-verbose test-coverage test-unit test-integration test-setup test-messages test-environment clean install lint security release release-dry-run release-push build-all commitlint commitlint-test

# Variáveis
GO=go
TEST_DIR=.
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Cores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m

# Ajuda
help: ## Mostrar esta ajuda
	@echo "🧪 Testes e Releases do Fin.Orion.Dev"
	@echo ""
	@echo "Comandos disponíveis:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# Instalação
install: ## Instalar dependências
	@echo "$(BLUE)[INFO]$(NC) Instalando dependências..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)[SUCCESS]$(NC) Dependências instaladas!"

# Testes
test: install ## Executar todos os testes
	@echo "$(BLUE)[INFO]$(NC) Executando todos os testes..."
	$(GO) test ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Todos os testes passaram!"

test-verbose: install ## Executar testes com verbose
	@echo "$(BLUE)[INFO]$(NC) Executando testes com verbose..."
	$(GO) test -v ./...
	@echo "$(GREEN)[SUCCESS]$(NC) Todos os testes passaram!"

test-coverage: install ## Executar testes com cobertura
	@echo "$(BLUE)[INFO]$(NC) Executando testes com cobertura..."
	$(GO) test -cover -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GO) tool cover -func=$(COVERAGE_FILE)
	$(GO) tool cover -html=$(COVERAGE_FILE) -o=$(COVERAGE_HTML)
	@echo "$(GREEN)[SUCCESS]$(NC) Relatório de cobertura gerado: $(COVERAGE_HTML)"

test-unit: install ## Executar apenas testes unitários
	@echo "$(BLUE)[INFO]$(NC) Executando testes unitários..."
	$(GO) test -v -run "Test.*" -short
	@echo "$(GREEN)[SUCCESS]$(NC) Testes unitários passaram!"

test-integration: install ## Executar apenas testes de integração
	@echo "$(BLUE)[INFO]$(NC) Executando testes de integração..."
	$(GO) test -v -run "Test.*Integration"
	@echo "$(GREEN)[SUCCESS]$(NC) Testes de integração passaram!"

test-setup: install ## Executar apenas testes de setup
	@echo "$(BLUE)[INFO]$(NC) Executando testes de setup..."
	$(GO) test -v -run "TestSetup"
	@echo "$(GREEN)[SUCCESS]$(NC) Testes de setup passaram!"

test-messages: install ## Executar apenas testes de mensagens
	@echo "$(BLUE)[INFO]$(NC) Executando testes de mensagens..."
	$(GO) test -v -run "TestMessages"
	@echo "$(GREEN)[SUCCESS]$(NC) Testes de mensagens passaram!"

test-environment: install ## Executar apenas testes de ambiente
	@echo "$(BLUE)[INFO]$(NC) Executando testes de ambiente..."
	$(GO) test -v -run "TestEnvironment"
	@echo "$(GREEN)[SUCCESS]$(NC) Testes de ambiente passaram!"

# Linting
lint: install ## Executar linting
	@echo "$(BLUE)[INFO]$(NC) Executando linting..."
	$(GO) vet ./...
	$(GO) fmt ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) golangci-lint não está instalado. Instalando..."; \
		GO111MODULE=on GOBIN=$(shell pwd)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		PATH=$(shell pwd)/bin:$$PATH golangci-lint run --timeout=5m; \
	fi
	@echo "$(GREEN)[SUCCESS]$(NC) Linting passou!"

# Segurança
security: install ## Executar verificação de segurança
	@echo "$(BLUE)[INFO]$(NC) Executando verificação de segurança..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) gosec não está instalado. Instale com: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Build
build-all: install ## Compilar para todas as plataformas
	@echo "$(BLUE)[INFO]$(NC) Compilando para todas as plataformas..."

	@echo "$(BLUE)[INFO]$(NC) Compilando para macOS Intel (amd64)..."
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-s -w" -o bin/orion-dev-darwin-amd64 cmd/main.go

	@echo "$(BLUE)[INFO]$(NC) Compilando para macOS Apple Silicon (arm64)..."
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="-s -w" -o bin/orion-dev-darwin-arm64 cmd/main.go

	@echo "$(BLUE)[INFO]$(NC) Compilando para Linux Intel (amd64)..."
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o bin/orion-dev-linux-amd64 cmd/main.go

	@echo "$(BLUE)[INFO]$(NC) Compilando para Linux ARM (arm64)..."
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags="-s -w" -o bin/orion-dev-linux-arm64 cmd/main.go

	@echo "$(BLUE)[INFO]$(NC) Compilando para Windows Intel (amd64)..."
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="-s -w" -o bin/orion-dev-windows-amd64.exe cmd/main.go

	@echo "$(BLUE)[INFO]$(NC) Compilando para Windows ARM (arm64)..."
	GOOS=windows GOARCH=arm64 $(GO) build -ldflags="-s -w" -o bin/orion-dev-windows-arm64.exe cmd/main.go

	@echo "$(GREEN)[SUCCESS]$(NC) Todos os binários foram compilados!"
	@echo "$(BLUE)[INFO]$(NC) Listando binários gerados:"
	@ls -lh bin/orion-dev-*

# Release
release: ## Criar release (requer versão como VERSION=1.2.3)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)[ERROR]$(NC) Versão é obrigatória. Use: make release VERSION=1.2.3"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Criando release $(VERSION)..."
	./scripts/release.sh -v $(VERSION)

release-dry-run: ## Testar release sem criar (requer versão como VERSION=1.2.3)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)[ERROR]$(NC) Versão é obrigatória. Use: make release-dry-run VERSION=1.2.3"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Testando release $(VERSION) (dry-run)..."
	./scripts/release.sh -v $(VERSION) -d

release-push: ## Criar release e fazer push (requer versão como VERSION=1.2.3)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)[ERROR]$(NC) Versão é obrigatória. Use: make release-push VERSION=1.2.3"; \
		exit 1; \
	fi
	@echo "$(BLUE)[INFO]$(NC) Criando release $(VERSION) e fazendo push..."
	./scripts/release.sh -v $(VERSION) -p

# Limpeza
clean: ## Limpar arquivos temporários
	@echo "$(BLUE)[INFO]$(NC) Limpando arquivos temporários..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f bin/test-*
	@echo "$(GREEN)[SUCCESS]$(NC) Limpeza concluída!"

clean-all: clean ## Limpar tudo incluindo binários
	@echo "$(BLUE)[INFO]$(NC) Limpando todos os binários..."
	rm -f bin/orion-dev-*
	@echo "$(GREEN)[SUCCESS]$(NC) Limpeza completa concluída!"

# Verificação de versão
version: ## Mostrar versão atual
	@echo "$(BLUE)[INFO]$(NC) Versão atual:"
	@git describe --tags --abbrev=0 2>/dev/null || echo "Nenhuma tag encontrada"

version-next: ## Sugerir próxima versão
	@echo "$(BLUE)[INFO]$(NC) Sugestão de próxima versão:"
	@if git describe --tags --abbrev=0 >/dev/null 2>&1; then \
		CURRENT=$$(git describe --tags --abbrev=0 | sed 's/v//'); \
		MAJOR=$$(echo $$CURRENT | cut -d. -f1); \
		MINOR=$$(echo $$CURRENT | cut -d. -f2); \
		PATCH=$$(echo $$CURRENT | cut -d. -f3); \
		echo "Versão atual: v$$CURRENT"; \
		echo "Sugestões:"; \
		echo "  Patch: v$$MAJOR.$$MINOR.$$((PATCH + 1))"; \
		echo "  Minor: v$$MAJOR.$$((MINOR + 1)).0"; \
		echo "  Major: v$$((MAJOR + 1)).0.0"; \
	else \
		echo "Primeira versão sugerida: v1.0.0"; \
	fi

# Commitlint
commitlint: install ## Validar último commit
	@echo "$(BLUE)[INFO]$(NC) Validando último commit..."
	@if [ -f "./bin/orion-dev" ]; then \
		./bin/orion-dev commitlint-last; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) Binário não encontrado. Compilando..."; \
		go build -o bin/orion-dev cmd/main.go; \
		./bin/orion-dev commitlint-last; \
	fi

commitlint-test: install ## Executar testes do commitlint
	@echo "$(BLUE)[INFO]$(NC) Executando testes do commitlint..."
	$(GO) test -v ./tests/commitlint_test.go
	@echo "$(GREEN)[SUCCESS]$(NC) Testes do commitlint passaram!"

commitlint-types: install ## Mostrar tipos válidos de commit
	@echo "$(BLUE)[INFO]$(NC) Tipos válidos de commit:"
	@if [ -f "./bin/orion-dev" ]; then \
		./bin/orion-dev commitlint-types; \
	else \
		echo "$(YELLOW)[WARNING]$(NC) Binário não encontrado. Compilando..."; \
		go build -o bin/orion-dev cmd/main.go; \
		./bin/orion-dev commitlint-types; \
	fi

install-hooks: install ## Instalar hooks do Git
	@echo "$(BLUE)[INFO]$(NC) Instalando hooks do Git..."
	@if [ -f "./scripts/install-hooks.sh" ]; then \
		./scripts/install-hooks.sh; \
	else \
		echo "$(RED)[ERROR]$(NC) Script de instalação não encontrado"; \
	fi