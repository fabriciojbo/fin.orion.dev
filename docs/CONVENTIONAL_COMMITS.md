# 📝 Conventional Commits e Commitlint

Este documento descreve o sistema de validação de commits implementado no projeto Fin.Orion.Dev, baseado no padrão **Conventional Commits**.

## 🎯 Visão Geral

O projeto utiliza um sistema nativo em Go para validar mensagens de commit, garantindo que sigam o padrão **Conventional Commits**. Isso permite:

- ✅ **Changelogs automáticos** baseados em commits
- 🔍 **Validação em tempo real** durante o desenvolvimento
- 📊 **Análise de mudanças** por tipo e escopo
- 🚀 **Releases automatizados** com changelog estruturado

## 📋 Padrão Conventional Commits

### Formato Básico

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Tipos Válidos

| Tipo       | Descrição                | Exemplo                                     |
| ---------- | ------------------------ | ------------------------------------------- |
| `feat`     | Nova funcionalidade      | `feat: adicionar autenticação JWT`          |
| `fix`      | Correção de bug          | `fix(auth): corrigir validação de senha`    |
| `docs`     | Documentação             | `docs: atualizar README`                    |
| `style`    | Formatação de código     | `style: formatar código com gofmt`          |
| `refactor` | Refatoração              | `refactor(api): reorganizar endpoints`      |
| `test`     | Testes                   | `test(service): adicionar testes unitários` |
| `chore`    | Manutenção               | `chore(deps): atualizar dependências`       |
| `perf`     | Melhorias de performance | `perf: otimizar consultas de banco`         |
| `ci`       | Integração contínua      | `ci: adicionar workflow de testes`          |
| `build`    | Build do sistema         | `build: configurar Docker multi-stage`      |
| `revert`   | Reverter commit          | `revert: reverter mudança problemática`     |

### Escopo (Opcional)

O escopo é opcional e deve ser escrito em minúsculas:

```
feat(auth): adicionar autenticação OAuth
fix(api): corrigir endpoint de usuários
test(service): adicionar testes para Service Bus
```

### Descrição

- **Obrigatória** e deve ser escrita em minúsculas
- **Máximo 72 caracteres** na primeira linha
- **Não deve terminar com ponto**
- **Use imperativo** ("adicionar" não "adicionado")

### Body (Opcional)

Para commits mais complexos, use o body para detalhar as mudanças:

```
feat(auth): adicionar autenticação JWT

Implementa autenticação JWT com refresh tokens.
Adiciona middleware de validação de token.
Inclui endpoints para login e logout.

- Adiciona package jwt
- Cria middleware auth
- Implementa refresh tokens
```

### Footer (Opcional)

Para referenciar issues, breaking changes, etc:

```
fix(auth): corrigir validação de token

Closes: #123
Fixes: #456
BREAKING CHANGE: Remove suporte a tokens antigos
```

### Breaking Changes

Para indicar mudanças que quebram compatibilidade:

```
feat!: remover suporte a API v1

BREAKING CHANGE: A API v1 foi removida. Use v2.
```

**Nota**: O `!` deve estar no final da mensagem, antes dos dois pontos da descrição.

## 🚀 Como Usar

### Comandos Disponíveis

#### 1. Validar Mensagem

```bash
# Validar mensagem específica
./bin/orion-dev commitlint "feat: adicionar nova funcionalidade"

# Validar mensagem interativamente
./bin/orion-dev commitlint
```

#### 2. Validar Último Commit

```bash
# Validar o último commit feito
./bin/orion-dev commitlint-last
```

#### 3. Ver Tipos Válidos

```bash
# Listar todos os tipos válidos
./bin/orion-dev commitlint-types
```

#### 4. Formatar Mensagem

```bash
# Formatar mensagem automaticamente
./bin/orion-dev commitlint-format feat auth "adicionar autenticação JWT"
# Resultado: feat(auth): adicionar autenticação JWT
```

#### 5. Hook de Validação

```bash
# Validar commit via hook (usado automaticamente)
./bin/orion-dev commitlint-hook
```

### Instalar Hooks do Git

Para validação automática em todos os commits:

```bash
# Instalar hooks automaticamente
./scripts/install-hooks.sh

# Forçar instalação (sobrescrever hooks existentes)
./scripts/install-hooks.sh -f

# Desinstalar hooks
./scripts/install-hooks.sh -u
```

## 🔧 Hooks Instalados

### commit-msg
Valida a mensagem de commit automaticamente:

```bash
git commit -m "feat: nova funcionalidade"
# ✅ Commit aceito

git commit -m "adicionar funcionalidade"
# ❌ Commit rejeitado - formato inválido
```

### pre-commit
Executa testes antes do commit:

```bash
git commit -m "feat: nova funcionalidade"
# 🧪 Executando testes...
# ✅ Testes passaram!
# ✅ Commit realizado
```

### pre-push
Valida o último commit antes do push:

```bash
git push origin main
# 🔍 Validando antes do push...
# ✅ Validação concluída!
# ✅ Push realizado
```

## 📝 Exemplos Práticos

### ✅ Mensagens Válidas

```bash
# Nova funcionalidade
feat: adicionar comando de proxy TLS

# Correção com escopo
fix(messages): corrigir validação de JSON

# Documentação
docs: atualizar README com novos comandos

# Testes com escopo
test(api): adicionar testes para endpoints

# Manutenção
chore(deps): atualizar dependências

# Breaking change
feat!: remover suporte a versão antiga

# Commit com body
feat(auth): implementar autenticação OAuth

Adiciona suporte completo a OAuth 2.0 com Google e GitHub.
Inclui refresh tokens e logout automático.

Closes: #123
```

### ❌ Mensagens Inválidas

```bash
# Tipo inválido
invalid: adicionar funcionalidade

# Sem descrição
feat:

# Descrição muito longa
feat: adicionar uma funcionalidade muito longa que excede o limite de caracteres permitido

# Termina com ponto
feat: adicionar funcionalidade.

# Formato incorreto
adicionar funcionalidade

# Escopo com maiúsculas
feat(Auth): adicionar autenticação
```

## ⚙️ Configuração

### Configuração Padrão

```go
config := &commitlint.Config{
    MaxSubjectLength: 72,    // Máximo de caracteres na descrição
    MinSubjectLength: 1,     // Mínimo de caracteres na descrição
    AllowedTypes:     ValidTypes, // Tipos permitidos
    RequireScope:     false, // Escopo obrigatório
    RequireBody:      false, // Body obrigatório
    RequireFooter:    false, // Footer obrigatório
}
```

### Configuração Customizada

```go
// Configuração mais rigorosa
config := &commitlint.Config{
    MaxSubjectLength: 50,
    MinSubjectLength: 5,
    AllowedTypes:     []commitlint.CommitType{
        commitlint.TypeFeat,
        commitlint.TypeFix,
        commitlint.TypeDocs,
    },
    RequireScope:     true,
    RequireBody:      true,
    RequireFooter:    false,
}
```

## 🧪 Testes

### Executar Testes do Commitlint

```bash
# Executar todos os testes
go test ./tests/commitlint_test.go

# Executar com verbose
go test -v ./tests/commitlint_test.go

# Executar com cobertura
go test -cover ./tests/commitlint_test.go
```

### Testes Disponíveis

- ✅ Validação de mensagens válidas
- ❌ Validação de mensagens inválidas
- 🔧 Formatação de mensagens
- 📋 Parse de diferentes formatos
- ⚙️ Configurações customizadas
- 🏷️ Tipos de commit válidos

## 🔄 Integração com Releases

O sistema de commitlint se integra perfeitamente com o sistema de releases:

### Changelog Automático

Commits válidos geram changelogs estruturados:

```markdown
## 🚀 Release 1.2.3

### 📝 Changes
- feat: adicionar comando de proxy TLS
- fix(messages): corrigir validação de JSON
- docs: atualizar README com novos comandos
- test(api): adicionar testes para endpoints

### 🔗 Issues & PRs
- Adicionar suporte a proxy TLS (#123)
- Corrigir validação de mensagens JSON (#124)
- Melhorar documentação (#125)
```

### Versionamento Semântico

Os tipos de commit influenciam o versionamento:

- `feat` → Incrementa versão MINOR
- `fix` → Incrementa versão PATCH
- `feat!` ou `BREAKING CHANGE` → Incrementa versão MAJOR

## 🚨 Troubleshooting

### Problemas Comuns

#### 1. Commit Rejeitado

```bash
git commit -m "adicionar funcionalidade"
# ❌ Commit rejeitado!

# Solução: Use o formato correto
git commit -m "feat: adicionar funcionalidade"
```

#### 2. Hook Não Funciona

```bash
# Verificar se hooks estão instalados
ls -la .git/hooks/

# Reinstalar hooks
./scripts/install-hooks.sh -f
```

#### 3. Validação Manual

```bash
# Validar mensagem antes do commit
./bin/orion-dev commitlint "feat: minha mensagem"

# Formatar mensagem
./bin/orion-dev commitlint-format feat scope "descrição"
```

### Logs de Debug

```bash
# Ver logs do hook
cat .git/hooks/commit-msg

# Testar hook manualmente
./bin/orion-dev commitlint-hook
```

## 📚 Recursos Adicionais

### Links Úteis

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Git Hooks](https://git-scm.com/docs/githooks)

### Comandos Úteis

```bash
# Ver tipos válidos
./bin/orion-dev commitlint-types

# Validar último commit
./bin/orion-dev commitlint-last

# Formatar mensagem
./bin/orion-dev commitlint-format feat auth "adicionar autenticação"

# Instalar hooks
./scripts/install-hooks.sh

# Executar testes
go test ./tests/commitlint_test.go
```

---

## 🎯 Próximos Passos

1. **Instale os hooks**: `./scripts/install-hooks.sh`
2. **Use Conventional Commits** em todos os commits
3. **Valide mensagens** antes de commitar
4. **Aproveite changelogs automáticos** nos releases

Para dúvidas ou problemas, consulte a documentação ou abra uma issue no repositório.