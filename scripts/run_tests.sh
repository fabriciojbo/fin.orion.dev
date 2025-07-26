#!/bin/bash

# 🧪 Script para executar testes do Fin.Orion.Dev
# Autor: Fin.Orion.Dev Team
# Data: $(date)

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir com cores
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Função para mostrar ajuda
show_help() {
    echo "🧪 Testes do Fin.Orion.Dev"
    echo ""
    echo "Uso: $0 [OPÇÕES]"
    echo ""
    echo "Opções:"
    echo "  -h, --help          Mostrar esta ajuda"
    echo "  -v, --verbose       Executar com verbose"
    echo "  -c, --coverage      Executar com cobertura"
    echo "  -u, --unit          Executar apenas testes unitários"
    echo "  -i, --integration   Executar apenas testes de integração"
    echo "  -s, --setup         Executar apenas testes de setup"
    echo "  -m, --messages      Executar apenas testes de mensagens"
    echo "  -e, --environment   Executar apenas testes de ambiente"
    echo "  -a, --all           Executar todos os testes (padrão)"
    echo ""
    echo "Exemplos:"
    echo "  $0                    # Executar todos os testes"
    echo "  $0 -v                 # Executar com verbose"
    echo "  $0 -c                 # Executar com cobertura"
    echo "  $0 -s                 # Executar apenas testes de setup"
    echo "  $0 -m -v              # Executar testes de mensagens com verbose"
}

# Variáveis
VERBOSE=false
COVERAGE=false
TEST_TYPE="all"

# Parse argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -u|--unit)
            TEST_TYPE="unit"
            shift
            ;;
        -i|--integration)
            TEST_TYPE="integration"
            shift
            ;;
        -s|--setup)
            TEST_TYPE="setup"
            shift
            ;;
        -m|--messages)
            TEST_TYPE="messages"
            shift
            ;;
        -e|--environment)
            TEST_TYPE="environment"
            shift
            ;;
        -a|--all)
            TEST_TYPE="all"
            shift
            ;;
        *)
            print_error "Opção desconhecida: $1"
            show_help
            exit 1
            ;;
    esac
done

# Verificar se estamos no diretório correto
if [[ ! -f "go.mod" ]]; then
    print_error "Este script deve ser executado na raiz do projeto"
    exit 1
fi

print_status "Iniciando execução dos testes..."

# Instalar dependências se necessário
print_status "Instalando dependências..."
go mod tidy
go mod download

# Construir comando de teste
TEST_CMD="go test"

if [[ "$VERBOSE" == "true" ]]; then
    TEST_CMD="$TEST_CMD -v"
fi

if [[ "$COVERAGE" == "true" ]]; then
    TEST_CMD="$TEST_CMD -cover"
    TEST_CMD="$TEST_CMD -coverprofile=coverage.out"
    TEST_CMD="$TEST_CMD -covermode=atomic"
fi

# Executar testes baseado no tipo
case $TEST_TYPE in
    "setup")
        print_status "Executando testes de setup..."
        $TEST_CMD -run "TestSetup"
        ;;
    "messages")
        print_status "Executando testes de mensagens..."
        $TEST_CMD -run "TestMessages"
        ;;
    "environment")
        print_status "Executando testes de ambiente..."
        $TEST_CMD -run "TestEnvironment"
        ;;
    "unit")
        print_status "Executando testes unitários..."
        $TEST_CMD -run "Test.*" -short
        ;;
    "integration")
        print_status "Executando testes de integração..."
        $TEST_CMD -run "Test.*Integration"
        ;;
    "all")
        print_status "Executando todos os testes..."
        $TEST_CMD ./...
        ;;
esac

# Verificar resultado
if [[ $? -eq 0 ]]; then
    print_success "Todos os testes passaram!"

        # Mostrar cobertura se solicitado
    if [[ "$COVERAGE" == "true" ]]; then
        if [[ -f "tests/coverage.out" ]]; then
            print_status "Gerando relatório de cobertura..."
            cd tests
            go tool cover -html=coverage.out -o coverage.html
            cd ..
            print_success "Relatório de cobertura gerado: tests/coverage.html"

            # Mostrar cobertura no terminal
            print_status "Cobertura de código:"
            cd tests
            go tool cover -func=coverage.out
            cd ..
        fi
    fi
else
    print_error "Alguns testes falharam!"
    exit 1
fi

print_success "Execução dos testes concluída!"