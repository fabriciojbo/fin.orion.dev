# 🚀 Guia de Releases - Fin.Orion.Dev

Este documento descreve o processo completo de criação de releases para o projeto Fin.Orion.Dev.

## 📋 Visão Geral

O sistema de releases é automatizado e baseado em **Semantic Versioning (SemVer)**. Quando uma tag é criada e enviada para o repositório, o GitHub Actions automaticamente:

1. ✅ Executa todos os testes
2. 🔨 Compila binários para 6 plataformas
3. 🔐 Gera checksums SHA256
4. 📝 Cria changelog baseado em commits
5. 🚀 Publica o release no GitHub

## 🏷️ Estratégia de Versionamento

### Semantic Versioning (SemVer)

Usamos o formato `MAJOR.MINOR.PATCH`:

- **MAJOR**: Mudanças incompatíveis com versões anteriores
- **MINOR**: Novas funcionalidades compatíveis
- **PATCH**: Correções de bugs compatíveis

**Exemplos:**
- `1.0.0` - Primeira versão estável
- `1.1.0` - Nova funcionalidade
- `1.1.1` - Correção de bug
- `2.0.0` - Breaking changes

### Conventional Commits

Para gerar changelogs automáticos, use commits no formato:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Tipos válidos:**
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Documentação
- `style`: Formatação
- `refactor`: Refatoração
- `test`: Testes
- `chore`: Manutenção

**Exemplos:**
```bash
git commit -m "feat: adicionar comando de proxy TLS"
git commit -m "fix(messages): corrigir validação de JSON"
git commit -m "docs: atualizar README com novos comandos"
```

## 🚀 Como Criar um Release

### Método 1: Script Automatizado (Recomendado)

```bash
# Criar tag localmente
./scripts/release.sh -v 1.2.3

# Criar tag e fazer push automaticamente
./scripts/release.sh -v 1.2.3 -p

# Com mensagem personalizada
./scripts/release.sh -v 1.2.3 -m "feat: nova funcionalidade de proxy" -p

# Dry-run (sem fazer alterações)
./scripts/release.sh -v 1.2.3 -d
```

### Método 2: Manual

```bash
# 1. Verificar se tudo está commitado
git status

# 2. Executar testes
go test -v ./...

# 3. Criar tag
git tag -a v1.2.3 -m "Release 1.2.3"

# 4. Fazer push da tag
git push origin v1.2.3
```

### Método 3: Via GitHub

1. Vá para **Releases** no GitHub
2. Clique em **"Create a new release"**
3. Crie uma nova tag (ex: `v1.2.3`)
4. Adicione descrição
5. Publique

## 📦 Plataformas Suportadas

Cada release inclui binários para:

| Plataforma              | Arquivo                       | Arquitetura |
| ----------------------- | ----------------------------- | ----------- |
| **macOS Intel**         | `orion-dev-darwin-amd64`      | x86_64      |
| **macOS Apple Silicon** | `orion-dev-darwin-arm64`      | ARM64       |
| **Linux Intel**         | `orion-dev-linux-amd64`       | x86_64      |
| **Linux ARM**           | `orion-dev-linux-arm64`       | ARM64       |
| **Windows Intel**       | `orion-dev-windows-amd64.exe` | x86_64      |
| **Windows ARM**         | `orion-dev-windows-arm64.exe` | ARM64       |

## 🔐 Verificação de Integridade

Cada binário inclui um arquivo `.sha256` para verificação:

```bash
# Download
wget https://github.com/user/repo/releases/download/v1.2.3/orion-dev-linux-amd64
wget https://github.com/user/repo/releases/download/v1.2.3/orion-dev-linux-amd64.sha256

# Verificar
sha256sum -c orion-dev-linux-amd64.sha256

# Tornar executável
chmod +x orion-dev-linux-amd64
```

## 📝 Changelog Automático

O changelog é gerado automaticamente baseado em:

1. **Conventional Commits** desde o último release
2. **Issues fechadas** desde o último release
3. **Pull Requests mergeados** desde o último release

### Exemplo de Changelog Gerado

```markdown
## 🚀 Release 1.2.3

### 📝 Changes
- feat: adicionar comando de proxy TLS
- fix(messages): corrigir validação de JSON
- docs: atualizar README com novos comandos
- test: adicionar testes para comandos de mensagens

### 🔗 Issues & PRs
- Adicionar suporte a proxy TLS (#123)
- Corrigir validação de mensagens JSON (#124)
- Melhorar documentação (#125)

### 📦 Downloads
- **macOS Intel (amd64)**: `orion-dev-darwin-amd64`
- **macOS Apple Silicon (arm64)**: `orion-dev-darwin-arm64`
- **Linux Intel (amd64)**: `orion-dev-linux-amd64`
- **Linux ARM (arm64)**: `orion-dev-linux-arm64`
- **Windows Intel (amd64)**: `orion-dev-windows-amd64.exe`
- **Windows ARM (arm64)**: `orion-dev-windows-arm64.exe`
```

## 🔧 Configuração do Workflow

### GitHub Actions

O workflow está configurado em `.github/workflows/release.yml`:

```yaml
on:
  push:
    tags:
      - 'v*' # Trigger em tags de versão

jobs:
  test:        # Executa testes
  build:       # Compila para 6 plataformas
  release:     # Cria o release
```

### Variáveis de Ambiente

- `GO_VERSION`: Versão do Go (1.24)
- `CGO_ENABLED`: Desabilitado para builds estáticos
- `GOOS/GOARCH`: Configurados para cada plataforma

## 🧪 Testes Antes do Release

Antes de cada release, são executados:

1. **Testes Unitários** - `go test -v ./...`
2. **Testes de Race** - `go test -race`
3. **Análise de Código** - `go vet ./...`
4. **Linting** - `golangci-lint`
5. **Build Test** - Testa compilação para todas as plataformas
6. **Security Scan** - `gosec` para vulnerabilidades

## 📊 Monitoramento

### Status do Workflow

- Acompanhe o progresso em **Actions** no GitHub
- Receba notificações de sucesso/falha
- Verifique logs detalhados de cada etapa

### Métricas

- **Tempo de Build**: ~5-10 minutos
- **Cobertura de Testes**: >80%
- **Tamanho dos Binários**: ~10-15MB cada

## 🚨 Troubleshooting

### Problemas Comuns

#### 1. Testes Falharam
```bash
# Executar testes localmente
go test -v ./...

# Verificar cobertura
go test -cover ./...
```

#### 2. Build Falhou
```bash
# Verificar dependências
go mod tidy
go mod download

# Testar build local
go build -o bin/test cmd/main.go
```

#### 3. Tag Já Existe
```bash
# Listar tags existentes
git tag -l

# Deletar tag local (se necessário)
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3
```

#### 4. Workflow Não Executou
- Verificar se a tag segue o padrão `v*`
- Verificar permissões do repositório
- Verificar configuração do GitHub Actions

### Logs de Debug

```bash
# Ver logs do workflow
# Vá para Actions > [Workflow] > [Job] > [Step]

# Ver logs detalhados
# Clique em "View logs" em cada etapa
```

## 📚 Recursos Adicionais

### Links Úteis

- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Go Build Constraints](https://golang.org/pkg/go/build/)

### Comandos Úteis

```bash
# Ver versão atual
git describe --tags --abbrev=0

# Ver commits desde último release
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Ver diferenças entre releases
git diff v1.2.2..v1.2.3

# Listar todas as tags
git tag -l --sort=-version:refname
```

---

## 🎯 Próximos Passos

1. **Configure o repositório** com as permissões necessárias
2. **Teste o workflow** criando uma tag de teste
3. **Documente mudanças** usando conventional commits
4. **Crie releases regulares** seguindo o processo

Para dúvidas ou problemas, abra uma issue no repositório.