#!/bin/bash

# 🔧 Script para instalar hooks do Git para Conventional Commits
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
    echo -e "${BLUE}🔧 Script para instalar hooks do Git${NC}"
    echo ""
    echo "Uso: $0 [OPÇÕES]"
    echo ""
    echo "Opções:"
    echo "  -h, --help               Mostrar esta ajuda"
    echo "  -f, --force              Forçar instalação (sobrescrever hooks existentes)"
    echo "  -u, --uninstall          Desinstalar hooks"
    echo ""
    echo "Exemplos:"
    echo "  $0                      # Instalar hooks"
    echo "  $0 -f                   # Forçar instalação"
    echo "  $0 -u                   # Desinstalar hooks"
    echo ""
}

# Função para verificar se estamos em um repositório git
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}❌ Erro: Não estamos em um repositório git${NC}"
        exit 1
    fi
}

# Função para obter o caminho do binário
get_binary_path() {
    # Tentar encontrar o binário compilado
    local binary_paths=(
        "./bin/orion-dev"
        "./bin/orion-dev.exe"
        "./orion-dev"
        "./orion-dev.exe"
    )

    for path in "${binary_paths[@]}"; do
        if [[ -f "$path" && -x "$path" ]]; then
            echo "$(pwd)/$path"
            return 0
        fi
    done

    # Se não encontrar, tentar compilar
    echo -e "${YELLOW}⚠️  Binário não encontrado. Tentando compilar...${NC}"
    if command -v go >/dev/null 2>&1; then
        go build -o bin/orion-dev cmd/main.go
        if [[ -f "./bin/orion-dev" ]]; then
            echo "$(pwd)/bin/orion-dev"
            return 0
        fi
    fi

    echo -e "${RED}❌ Erro: Não foi possível encontrar ou compilar o binário${NC}"
    exit 1
}

# Função para instalar hooks
install_hooks() {
    local force=$1
    local binary_path=$(get_binary_path)
    local hooks_dir=".git/hooks"

    echo -e "${BLUE}🔧 Instalando hooks do Git...${NC}"
    echo ""

    # Criar diretório de hooks se não existir
    mkdir -p "$hooks_dir"

    # Hook commit-msg
    local commit_msg_hook="$hooks_dir/commit-msg"
    if [[ -f "$commit_msg_hook" && "$force" != "true" ]]; then
        echo -e "${YELLOW}⚠️  Hook commit-msg já existe. Use -f para sobrescrever.${NC}"
    else
        cat > "$commit_msg_hook" << EOF
#!/bin/bash
# Hook para validar mensagens de commit usando Conventional Commits
# Instalado automaticamente pelo Fin.Orion.Dev

set -e

# Executar validador de commit
if ! "$binary_path" commitlint-hook; then
    echo ""
    echo "💡 Dica: Use o comando para formatar sua mensagem:"
    echo "   $binary_path commitlint-format <type> <scope> <description>"
    echo ""
    echo "📋 Tipos válidos:"
    echo "   $binary_path commitlint-types"
    echo ""
    exit 1
fi
EOF
        chmod +x "$commit_msg_hook"
        echo -e "${GREEN}✅ Hook commit-msg instalado${NC}"
    fi

    # Hook pre-commit (opcional)
    local pre_commit_hook="$hooks_dir/pre-commit"
    if [[ -f "$pre_commit_hook" && "$force" != "true" ]]; then
        echo -e "${YELLOW}⚠️  Hook pre-commit já existe. Use -f para sobrescrever.${NC}"
    else
        cat > "$pre_commit_hook" << EOF
#!/bin/bash
# Hook para executar testes antes do commit
# Instalado automaticamente pelo Fin.Orion.Dev

set -e

echo "🧪 Executando testes antes do commit..."

# Executar testes
if command -v go >/dev/null 2>&1; then
    if ! go test ./...; then
        echo ""
        echo "❌ Testes falharam. Commit cancelado."
        echo "💡 Corrija os testes antes de fazer commit."
        exit 1
    fi
    echo "✅ Testes passaram!"
else
    echo "⚠️  Go não encontrado. Pulando testes."
fi

# Verificar formatação
if command -v go >/dev/null 2>&1; then
    echo "🔍 Verificando formatação..."
    if ! go fmt ./...; then
        echo "⚠️  Problemas de formatação encontrados."
    fi
fi
EOF
        chmod +x "$pre_commit_hook"
        echo -e "${GREEN}✅ Hook pre-commit instalado${NC}"
    fi

    # Hook pre-push (opcional)
    local pre_push_hook="$hooks_dir/pre-push"
    if [[ -f "$pre_push_hook" && "$force" != "true" ]]; then
        echo -e "${YELLOW}⚠️  Hook pre-push já existe. Use -f para sobrescrever.${NC}"
    else
        cat > "$pre_push_hook" << EOF
#!/bin/bash
# Hook para validar antes do push
# Instalado automaticamente pelo Fin.Orion.Dev

set -e

echo "🔍 Validando antes do push..."

# Validar último commit
if ! "$binary_path" commitlint-last; then
    echo ""
    echo "❌ Último commit é inválido. Push cancelado."
    echo "💡 Corrija a mensagem do commit antes de fazer push."
    exit 1
fi

echo "✅ Validação concluída!"
EOF
        chmod +x "$pre_push_hook"
        echo -e "${GREEN}✅ Hook pre-push instalado${NC}"
    fi

    echo ""
    echo -e "${GREEN}🎉 Hooks instalados com sucesso!${NC}"
    echo ""
    echo -e "${BLUE}📋 Hooks instalados:${NC}"
    echo "  • commit-msg  - Valida mensagens de commit"
    echo "  • pre-commit  - Executa testes antes do commit"
    echo "  • pre-push    - Valida antes do push"
    echo ""
    echo -e "${BLUE}💡 Comandos úteis:${NC}"
    echo "  $binary_path commitlint-types     # Ver tipos válidos"
    echo "  $binary_path commitlint-format    # Formatar mensagem"
    echo "  $binary_path commitlint-last      # Validar último commit"
    echo ""
}

# Função para desinstalar hooks
uninstall_hooks() {
    local hooks_dir=".git/hooks"
    local hooks=("commit-msg" "pre-commit" "pre-push")

    echo -e "${BLUE}🗑️  Desinstalando hooks do Git...${NC}"
    echo ""

    for hook in "${hooks[@]}"; do
        local hook_path="$hooks_dir/$hook"
        if [[ -f "$hook_path" ]]; then
            # Verificar se é nosso hook (contém "Fin.Orion.Dev")
            if grep -q "Fin.Orion.Dev" "$hook_path"; then
                rm "$hook_path"
                echo -e "${GREEN}✅ Hook $hook removido${NC}"
            else
                echo -e "${YELLOW}⚠️  Hook $hook não foi instalado por este script${NC}"
            fi
        else
            echo -e "${YELLOW}⚠️  Hook $hook não existe${NC}"
        fi
    done

    echo ""
    echo -e "${GREEN}🎉 Hooks desinstalados com sucesso!${NC}"
}

# Variáveis
FORCE=false
UNINSTALL=false

# Parse argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -u|--uninstall)
            UNINSTALL=true
            shift
            ;;
        *)
            echo -e "${RED}❌ Opção desconhecida: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Verificar se estamos em um repositório git
check_git_repo

# Executar ação
if [[ "$UNINSTALL" == true ]]; then
    uninstall_hooks
else
    install_hooks "$FORCE"
fi