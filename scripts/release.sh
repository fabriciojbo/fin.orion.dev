#!/bin/bash

# 🚀 Script de Release para Fin.Orion.Dev
# Autor: Fin.Orion.Dev Team

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Função de ajuda
show_help() {
    echo -e "${BLUE}🚀 Script de Release para Fin.Orion.Dev${NC}"
    echo ""
    echo "Uso: $0 [OPÇÕES]"
    echo ""
    echo "Opções:"
    echo "  -v, --version VERSION    Versão para release (ex: 1.2.3)"
    echo "  -m, --message MESSAGE    Mensagem do commit (opcional)"
    echo "  -p, --push               Fazer push da tag automaticamente"
    echo "  -d, --dry-run            Executar sem fazer alterações"
    echo "  -h, --help               Mostrar esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0 -v 1.2.3              Criar tag v1.2.3"
    echo "  $0 -v 1.2.3 -p           Criar tag e fazer push"
    echo "  $0 -v 1.2.3 -m 'feat: nova funcionalidade' -p"
    echo ""
}

# Função para validar versão semântica
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo -e "${RED}❌ Erro: Versão deve seguir o formato MAJOR.MINOR.PATCH (ex: 1.2.3)${NC}"
        exit 1
    fi
}

# Função para verificar se há mudanças não commitadas
check_uncommitted_changes() {
    if [[ -n $(git status --porcelain) ]]; then
        echo -e "${YELLOW}⚠️  Aviso: Há mudanças não commitadas${NC}"
        git status --short
        echo ""
        read -p "Deseja continuar mesmo assim? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${RED}❌ Release cancelado${NC}"
            exit 1
        fi
    fi
}

# Função para verificar se a tag já existe
check_existing_tag() {
    local version=$1
    local tag="v$version"

    if git tag -l | grep -q "^$tag$"; then
        echo -e "${RED}❌ Erro: Tag $tag já existe${NC}"
        exit 1
    fi
}

# Função para verificar testes
run_tests() {
    echo -e "${BLUE}🧪 Executando testes...${NC}"

    if ! go test -v ./...; then
        echo -e "${RED}❌ Testes falharam. Abortando release.${NC}"
        exit 1
    fi

    echo -e "${GREEN}✅ Testes passaram${NC}"
}

# Função para verificar build
test_build() {
    echo -e "${BLUE}🔨 Testando build...${NC}"

    # Testar build para plataforma atual
    if ! go build -o bin/test-release cmd/main.go; then
        echo -e "${RED}❌ Build falhou. Abortando release.${NC}"
        exit 1
    fi

    # Limpar arquivo de teste
    rm -f bin/test-release

    echo -e "${GREEN}✅ Build testado com sucesso${NC}"
}

# Função para criar tag
create_tag() {
    local version=$1
    local message=$2
    local tag="v$version"

    echo -e "${BLUE}🏷️  Criando tag $tag...${NC}"

    if [[ -n "$message" ]]; then
        git tag -a "$tag" -m "$message"
    else
        git tag -a "$tag" -m "Release $tag"
    fi

    echo -e "${GREEN}✅ Tag $tag criada localmente${NC}"
}

# Função para fazer push
push_tag() {
    local version=$1
    local tag="v$version"

    echo -e "${BLUE}📤 Fazendo push da tag $tag...${NC}"

    if git push origin "$tag"; then
        echo -e "${GREEN}✅ Tag $tag enviada para o repositório remoto${NC}"
        echo -e "${GREEN}🚀 GitHub Actions irá criar o release automaticamente${NC}"
    else
        echo -e "${RED}❌ Erro ao fazer push da tag${NC}"
        exit 1
    fi
}

# Função para mostrar informações do release
show_release_info() {
    local version=$1
    local tag="v$version"

    echo ""
    echo -e "${GREEN}🎉 Release $tag preparado com sucesso!${NC}"
    echo ""
    echo -e "${BLUE}📋 Próximos passos:${NC}"
    echo "1. GitHub Actions irá executar automaticamente"
    echo "2. Binários serão compilados para todas as plataformas"
    echo "3. Release será criado em: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/releases"
    echo ""
    echo -e "${BLUE}📦 Binários que serão gerados:${NC}"
    echo "- orion-dev-darwin-amd64 (macOS Intel)"
    echo "- orion-dev-darwin-arm64 (macOS Apple Silicon)"
    echo "- orion-dev-linux-amd64 (Linux Intel)"
    echo "- orion-dev-linux-arm64 (Linux ARM)"
    echo "- orion-dev-windows-amd64.exe (Windows Intel)"
    echo "- orion-dev-windows-arm64.exe (Windows ARM)"
    echo ""
}

# Variáveis
VERSION=""
MESSAGE=""
PUSH=false
DRY_RUN=false

# Parse argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -m|--message)
            MESSAGE="$2"
            shift 2
            ;;
        -p|--push)
            PUSH=true
            shift
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Opção desconhecida: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Verificar se versão foi fornecida
if [[ -z "$VERSION" ]]; then
    echo -e "${RED}❌ Erro: Versão é obrigatória${NC}"
    show_help
    exit 1
fi

# Validar versão
validate_version "$VERSION"

# Verificar se estamos no diretório correto
if [[ ! -f "go.mod" ]] || [[ ! -f "cmd/main.go" ]]; then
    echo -e "${RED}❌ Erro: Execute este script no diretório raiz do projeto${NC}"
    exit 1
fi

# Verificar se git está disponível
if ! command -v git &> /dev/null; then
    echo -e "${RED}❌ Erro: Git não está instalado${NC}"
    exit 1
fi

# Verificar se estamos em um repositório git
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ Erro: Não estamos em um repositório git${NC}"
    exit 1
fi

echo -e "${BLUE}🚀 Iniciando processo de release para versão $VERSION${NC}"
echo ""

# Verificar mudanças não commitadas
check_uncommitted_changes

# Verificar se tag já existe
check_existing_tag "$VERSION"

# Executar testes
run_tests

# Testar build
test_build

# Se for dry-run, parar aqui
if [[ "$DRY_RUN" == true ]]; then
    echo -e "${YELLOW}🔍 Dry-run: Nenhuma alteração foi feita${NC}"
    echo -e "${GREEN}✅ Release $VERSION está pronto para ser criado${NC}"
    exit 0
fi

# Criar tag
create_tag "$VERSION" "$MESSAGE"

# Fazer push se solicitado
if [[ "$PUSH" == true ]]; then
    push_tag "$VERSION"
fi

# Mostrar informações finais
show_release_info "$VERSION"