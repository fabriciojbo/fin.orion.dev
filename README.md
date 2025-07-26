# 🚀 Ambiente de Desenvolvimento Orion

> Ambiente completo de desenvolvimento e testes para o sistema **Orion Functions** e **Orion API**, simulando toda a infraestrutura Azure necessária localmente com CLI nativo em Go.

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)](https://docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/fabriciojbo/fin.orion.dev)](https://github.com/fabriciojbo/fin.orion.dev/releases)
[![CI](https://img.shields.io/github/actions/workflow/status/fabriciojbo/fin.orion.dev/ci.yml?branch=main)](https://github.com/fabriciojbo/fin.orion.dev/actions)

## 📋 Índice

- [🎯 Visão Geral](#-visão-geral)
- [🏗️ Arquitetura](#️-arquitetura)
- [⚡️ Início Rápido](#-início-rápido)
- [🔧 Configuração](#-configuração)
- [🚀 Comandos](#-comandos)
- [📨 Mensagens e Filas](#-mensagens-e-filas)
- [🧪 Testes](#-testes)
- [📊 Monitoramento](#-monitoramento)
- [🛠️ Desenvolvimento](#️-desenvolvimento)
- [🗄️ Banco de Dados](#️-banco-de-dados)
- [🔍 Troubleshooting](#-troubleshooting)
- [📚 Documentação](#-documentação)
- [🏷️ Releases e Versionamento](#️-releases-e-versionamento)
- [📝 Conventional Commits](#-conventional-commits)

---

## 🎯 Visão Geral

Este ambiente simula completamente o ecossistema Orion com todos os serviços necessários para desenvolvimento local:

### 🐳 Serviços Disponíveis

| Serviço                  | Porta   | Descrição                    |
| ------------------------ | ------- | ---------------------------- |
| **Orion Functions**      | `7071`  | Azure Functions local        |
| **Orion API**            | `3333`  | API REST principal           |
| **Service Bus Emulator** | `5672`  | Filas e tópicos de mensagens |
| **Azurite Storage**      | `10000` | Azure Storage local          |
| **PostgreSQL**           | `5432`  | Banco de dados principal     |
| **SQL Server Edge**      | `1433`  | Banco de dados secundário    |

---

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Fin.Orion.Dev - Ambiente Completo                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │   Orion         │  │   Orion API     │  │   Service Bus   │              │
│  │   Functions     │◄─┤   :3333         │◄─┤   Emulator      │              │
│  │   :7071         │  │                 │  │   :5672         │              │
│  │                 │  │                 │  │                 │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│           │                     │                     │                     │
│           │                     │                     │                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │   Azurite       │  │   PostgreSQL    │  │   SQL Server    │              │
│  │   Storage       │  │   :5432         │  │   Edge          │              │
│  │   :10000        │  │                 │  │   :1433         │              │
│  │                 │  │                 │  │                 │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ⚡️ Início Rápido

### 📋 Pré-requisitos

```bash
# Dependências obrigatórias
docker --version          # Docker 20.10+
docker-compose --version  # Docker Compose 2.0+
go --version              # Go 1.24+
```

### 🚀 Setup Inicial

```bash
# 1. Compilar aplicação Go
go build -o bin/orion-dev cmd/main.go

# 2. Configurar ambiente (primeira vez)
./bin/orion-dev setup

# 3. Iniciar todos os serviços
./bin/orion-dev start

# 4. Verificar status
./bin/orion-dev status

# 5. Testar funcionalidade
./bin/orion-dev check-messages
```
---

## 🔧 Configuração

### 📁 Estrutura do Projeto

```
Fin.Orion.Dev/
├── 📁 bin/                           # Binário da aplicação
│   ├── .gitkeep                      # Mantém a pasta vazia
│   └── orion-dev                     # Binário da aplicação
├── 📁 docs/                          # Documentação
│   ├── CONVENTIONAL_COMMITS.md       # Conventional Commits
│   ├── RELEASES.md                   # Releases
│   └── TESTS.md                      # Testes
├── 📁 cmd/                           # Ponto de entrada da aplicação
│   └── main.go                       # Arquivo principal
├── 📁 internal/                      # Código interno da aplicação
│   ├── commands/                     # Comandos CLI
│   │   ├── root.go                   # Comando raiz
│   │   ├── setup.go                  # Comando setup
│   │   ├── start.go                  # Comando start
│   │   ├── stop.go                   # Comando stop
│   │   ├── status.go                 # Comando status
│   │   └── messages.go               # Comandos de mensagens
│   ├── commitlint/                   # Commitlint
│   │   └── validator.go              # Validador de commits
│   ├── proxy/                        # Proxy Service Bus
│   │   └── servicebus-proxy.go       # Proxy Service Bus
│   ├── servicebus/                   # Service Bus
│   │   └── client.go                 # Cliente Azure Service Bus
│   └── utils/                        # Utilitários
│       └── network.go                # Funções de rede
├── 📁 docker/                        # Configurações Docker
│   ├── container/
│   │   ├── Dockerfile.api            # Orion API
│   │   ├── Dockerfile.functions      # Orion Functions
│   │   └── .dockerignore             # Ignorar arquivos Docker
│   ├── database/
│   │   ├── certs/                    # Certificados PostgreSQL
│   │   │   ├── .gitkeep              # Mantém a pasta vazia
│   │   │   ├── server.crt            # Certificado PostgreSQL
│   │   │   └── server.key            # Chave privada PostgreSQL
│   │   ├── init-postgres.sql         # Script de inicialização
│   │   └── postgres.conf             # Configuração PostgreSQL
│   ├── service-bus/
│   │   ├── certs/                    # Certificados proxy Service Bus
│   │   │   ├── .gitkeep              # Mantém a pasta vazia
│   │   │   ├── servicebus-proxy.crt  # Certificado proxy Service Bus
│   │   │   └── servicebus-proxy.key  # Chave privada proxy Service Bus
│   │   └── config.json               # Configuração Service Bus
├── 📁 messages/                      # Arquivos JSON de teste
│   ├── rec_payment_order_fail.json   # Mensagem de exemplo
│   └── .gitkeep                      # Mantém a pasta vazia
├── 📁 scripts/                       # Scripts shell
│   ├── install-hooks.sh              # Instalação dos hooks
│   ├── release.sh                    # Release
│   └── tests.sh                      # Testes unitários
├── 📁 tests/                         # Testes unitários
├── 📄 .editorconfig                  # EditorConfig
├── 📄 .env                           # Variáveis de ambiente
├── 📄 .env.example                   # Exemplo de variáveis de ambiente
├── 📄 LICENSE                        # Licença
├── 📄 docker-compose.yml             # Orquestração Docker
├── 📄 go.mod                         # Dependências Go
├── 📄 go.sum                         # Checksums das dependências
├── 📄 local.settings.json            # Configuração Orion Functions
├── 📄 Makefile                       # Makefile
├── 📄 README.md                      # README
└── 📄 .gitignore                     # Arquivos ignorados
```

### 🔐 Variáveis de Ambiente

#### Orion API (`.env`)

```bash
# Configurações básicas
PORT=3333
ENV=HMG
API_KEY=FAKE-API-KEY

# JWT
JWT_SECRET=FAKE-JWT-SECRET
JWT_EXPIRE_MS=3600

# Application Insights
APPLICATIONINSIGHTS_CONNECTION_STRING="InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://brazilsouth-0.in.applicationinsights.azure.com/;LiveEndpoint=https://brazilsouth.livediagnostics.monitor.azure.com/;ApplicationId=00000000-0000-0000-0000-000000000000"

# PostgreSQL
PG_HOST="orion-database"
PG_PORT=5432
PG_USERNAME="postgres"
PG_PASSWORD="postgres"
PG_DATABASE="orion_database"
PG_SCHEMA="orionlocal"

# Pismo
PISMO_URL="https://sandbox.pismolabs.io"
PISMO_SERVER_KEY="YOUR_PISMO_SERVER_KEY"
PISMO_SERVER_SECRET="YOUR_PISMO_SERVER_SECRET"
PISMO_PROGRAM_ID="YOUR_PISMO_PROGRAM_ID"

# PIX
PIX_QRCODE_URL="pix-qrcode-h.magfinancas.com.br"

# Key Vault
KV_RESOURCE_NAME="YOUR_KV_RESOURCE_NAME"
KV_PISMOCERT_NAME="YOUR_KV_PISMOCERT_NAME"

# Contas
PAYMENT_ACCOUNT_SEQ=000000000000
BANK_BRANCH=0000

# MAG IP
MAG_IP_PAYMENT_ACCOUNT_ID="00000000-0000-0000-0000-000000000000"
MAG_IP_PISMO_ACCOUNT_ID=000000000000
MAG_IP_ACCOUNT_BRANCH=0000
MAG_IP_ACCOUNT_NUMBER=00000000

# Cron
DISABLE_CRON=false
CONSULT_BILLET_CRON="0 */1 * * *"

# Celcoin
CELCOIN_URL="YOUR_CELCOIN_URL"
CELCOIN_CLIENT_ID="YOUR_CELCOIN_CLIENT_ID"
CELCOIN_CLIENT_SECRET="YOUR_CELCOIN_CLIENT_SECRET"

# Service Bus
SB_CNT_STR="Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=FAKE-SAS-KEY-VALUE;UseDevelopmentEmulator=true"
SB_QUEUE_NAME="sbq.pismo.onboarding.succeeded"

# Azure Storage
AZURE_STORAGE_BLOB_STRING_CONNECTION="DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://azure-storage:10000/devstoreaccount1;QueueEndpoint=http://azure-storage:10001/devstoreaccount1;TableEndpoint=http://azure-storage:10002/devstoreaccount1;"
AZURE_STORAGE_BLOB_ACCOUNT_UPLOADS_CONTAINER="account-uploads"

# RSA Keys
RSA_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\nYOUR_RSA_PUBLIC_KEY\n-----END PUBLIC KEY-----"
RSA_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYOUR_RSA_PRIVATE_KEY\n-----END PRIVATE KEY-----"

# Magfinancas
MAGFINANCAS_SENSEDIA_BASE_URL="YOUR_MAGFINANCAS_SENSEDIA_BASE_URL"
MAGFINANCAS_SENSEDIA_BEARER_TOKEN="YOUR_MAGFINANCAS_SENSEDIA_BEARER_TOKEN"

# SQL Server Edge & Service Bus Emulator
ACCEPT_EULA="Y"
MSSQL_SA_PASSWORD="YOUR_MSSQL_SA_PASSWORD"
```

#### Orion Functions (`local.settings.json`)

```json
{
  "IsEncrypted": false,
  "Values": {
    "DEBUG": 1,
    "AzureWebJobsFeatureFlags": "EnableWorkerIndexing",
    "APPLICATIONINSIGHTS_CONNECTION_STRING": "InstrumentationKey=00000000-0000-0000-0000-000000000000;IngestionEndpoint=https://brazilsouth-0.in.applicationinsights.azure.com/;LiveEndpoint=https://brazilsouth.livediagnostics.monitor.azure.com/",
    "SB_CONN_STR": "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true",
    "SB_INPUT_QUEUE": "sbq.pismo.transaction.creation",
    "SB_OUTPUT_TOPIC": "sbt.orion.core",
    "PG_HOST": "orion-database",
    "PG_PORT": "5432",
    "PG_USERNAME": "postgres",
    "PG_PASSWORD": "postgres",
    "PG_DATABASE": "orion_database",
    "PG_SCHEMA": "orionlocal,orionlocal,orionlocal,orionlocal",
    "PG_SCHEMA_PIX": "orionlocal",
    "PISMO_URL": "https://sandbox.pismolabs.io",
    "PISMO_SERVER_KEY": "YOUR_PISMO_SERVER_KEY",
    "PISMO_SERVER_SECRET": "YOUR_PISMO_SERVER_SECRET",
    "PISMO_PROGRAM_ID": "YOUR_PISMO_PROGRAM_ID",
    "ORION_URL": "http://orion-api:3333",
    "ORION_API_KEY": "FAKE-API-KEY",
    "PIX_STG_CONN_STR": "DefaultEndpointsProtocol=https;AccountName=YOUR_ACCOUNT_NAME;AccountKey=YOUR_ACCOUNT_KEY;EndpointSuffix=core.windows.net"
  }
}
```

### 📨 Filas e Tópicos Configurados

| Fila/Tópico                                | Descrição                     | Tipo   |
| ------------------------------------------ | ----------------------------- | ------ |
| `sbq.pismo.onboarding.succeeded`           | Onboarding bem-sucedido       | Fila   |
| `sbq.pismo.transaction.creation`           | Criação de transações         | Fila   |
| `sbq.pismo.pix.transaction.in`             | Transações PIX IN             | Fila   |
| `sbq.pismo.all`                            | Todas as mensagens Pismo      | Fila   |
| `sbq.orion.pixqrcode.persist`              | Persistência QR Codes         | Fila   |
| `sbq.orion.transaction.chained`            | Transações encadeadas         | Fila   |
| `sbq.orion.billet-payment.verify`          | Verificação boletos           | Fila   |
| `sbq.pismo.authorization.cancelation`      | Autorizações canceladas       | Fila   |
| `sbq.pismo.ted.transaction`                | Transações TED                | Fila   |
| `sbq.pix.recurrence.payment.order.failure` | Falhas pagamentos recorrentes | Fila   |
| `sbt.orion.core`                           | Tópico principal Orion        | Tópico |

---

## 🚀 Comandos

### 📋 Comandos CLI Go (Principais)

```bash
# =============================================================================
# COMANDOS PRINCIPAIS
# =============================================================================

./bin/orion-dev setup          # Configurar ambiente inicial
./bin/orion-dev start          # Iniciar ambiente completo
./bin/orion-dev stop           # Parar ambiente
./bin/orion-dev status         # Ver status dos containers
./bin/orion-dev list           # Listar recursos disponíveis

# =============================================================================
# PROXY SERVICE BUS (OBRIGATÓRIO PARA MENSAGENS)
# =============================================================================

./bin/orion-dev proxy          # Iniciar proxy TLS (5671 -> 5672)

# =============================================================================
# COMANDOS DE MENSAGENS
# =============================================================================

./bin/orion-dev push-message <fila> <arquivo>  # Enviar mensagem para fila
./bin/orion-dev check-messages                 # Verificar mensagens do Service Bus
./bin/orion-dev check-topic [subscription]     # Verificar mensagens do tópico
./bin/orion-dev check-queue <fila>             # Verificar mensagens da fila

# =============================================================================
# COMANDOS DE JSON
# =============================================================================

./bin/orion-dev validate-json <arquivo>        # Validar arquivo JSON
./bin/orion-dev format-json <arquivo>          # Formatar arquivo JSON
./bin/orion-dev show-json <arquivo>            # Mostrar JSON formatado

# =============================================================================
# COMANDOS DE LIMPEZA
# =============================================================================

./bin/orion-dev stop --clean   # Parar ambiente e limpar recursos

# =============================================================================
# AJUDA
# =============================================================================

./bin/orion-dev help           # Ver todos os comandos
./bin/orion-dev <comando> --help  # Ajuda específica do comando
```

### 📋 Comandos Docker Compose (Auxiliares)

```bash
# =============================================================================
# COMANDOS DE CONTAINERS
# =============================================================================

# Ver logs dos containers
./bin/orion-dev logs

# Ver logs específicos
./bin/orion-dev logs --service orion-functions
./bin/orion-dev logs --service orion-api
./bin/orion-dev logs --service emulator

# Reconstruir containers
./bin/orion-dev build

# Reconstruir específico
./bin/orion-dev rebuild-functions
./bin/orion-dev rebuild-api

# Acessar shell dos containers
./bin/orion-dev shell --service orion-functions
./bin/orion-dev shell --service orion-api

# =============================================================================
# COMANDOS DE DESENVOLVIMENTO
# =============================================================================

# Verificar conectividade dos serviços
curl http://localhost:7071  # Orion Functions
curl http://localhost:3333  # Orion API
curl http://localhost:10000 # Azurite Storage

# Testar endpoints específicos
curl http://localhost:7071/cob/test-id
curl http://localhost:7071/cobv/test-id
curl -H 'X-API-Key: FAKE-API-KEY' http://localhost:3333/health
```

### 📋 Comandos Go (Desenvolvimento)

```bash
# =============================================================================
# COMPILAÇÃO PARA DIFERENTES PLATAFORMAS
# =============================================================================

# Compilar o binário para Unix (Mac Intel) (Architecture: amd64)
GOOS=darwin GOARCH=amd64 go build -o bin/orion-dev cmd/main.go

# Compilar o binário para Unix (Mac Apple Silicon) (Architecture: arm64)
GOOS=darwin GOARCH=arm64 go build -o bin/orion-dev cmd/main.go

# Compilar o binário para Unix (Linux) (Architecture: amd64)
GOOS=linux GOARCH=amd64 go build -o bin/orion-dev cmd/main.go

# Compilar o binário para Windows (x64) (Architecture: amd64)
GOOS=windows GOARCH=amd64 go build -o bin/orion-dev.exe cmd/main.go

# =============================================================================
# COMANDOS DE CONFIGURAÇÃO
# =============================================================================

# Executar comandos de configuração
./bin/orion-dev setup
./bin/orion-dev start
./bin/orion-dev status

# =============================================================================
# COMANDOS DE MENSAGENS
# =============================================================================

# Executar comandos de mensagens
./bin/orion-dev push-message <fila> <arquivo>
./bin/orion-dev check-messages
./bin/orion-dev check-queue <fila>
./bin/orion-dev check-topic [subscription]

# =============================================================================
# AJUDA
# =============================================================================

# Executar com ajuda
./bin/orion-dev help
./bin/orion-dev <comando> --help

# =============================================================================
# COMANDOS DE TESTE
# =============================================================================

# Executar testes
go test ./...

# Executar testes específicos
go test ./internal/commands/...
go test ./internal/servicebus/...
go test ./internal/utils/...

# Executar com cobertura
go test -cover ./...
```

---

## 📨 Mensagens e Filas

### 📤 Enviar Mensagens

```bash
# 1. Criar arquivo JSON de teste
echo '{"test": "message", "data": "example"}' > messages/test.json

# 2. Enviar para fila
./bin/orion-dev push-message sbq.pismo.transaction.creation test.json

# 3. Verificar mensagens
./bin/orion-dev check-queue sbq.pismo.transaction.creation
```

### 📥 Verificar Mensagens

```bash
# Verificar tópico
./bin/orion-dev check-topic

# Verificar fila específica
./bin/orion-dev check-queue sbq.orion.transaction.chained

# Listar recursos disponíveis
./bin/orion-dev list

# Verificar status geral
./bin/orion-dev check-messages
```

---

## 🧪 Testes

### 🚀 Testes Rápidos

```bash
# Teste completo do ambiente
./bin/orion-dev check-messages

# Verificar saúde dos serviços
./bin/orion-dev status

# Testar Orion Functions
curl http://localhost:7071/cob/test-id
curl http://localhost:7071/cobv/test-id

# Testar Orion API
curl http://localhost:3333
```

### 📨 Testes de Mensagens

```bash
# 1. Enviar mensagem de teste
./bin/orion-dev push-message sbq.pismo.transaction.creation test.json

# 2. Verificar se foi processada
./bin/orion-dev check-queue sbq.pismo.transaction.creation

# 3. Ver logs
./bin/orion-dev logs
```

### 🔗 Testes de Integração

```bash
# Teste PIX Recurrence completo
./bin/orion-dev push-message sbq.pix.recurrence.payment.order.failure pix-recurrence.json
sleep 5
./bin/orion-dev check-queue sbq.pix.recurrence.payment.order.failure

# Teste transação encadeada
./bin/orion-dev push-message sbq.orion.transaction.chained transaction-chained.json
sleep 5
./bin/orion-dev check-queue sbq.orion.transaction.chained
```

---

## 📊 Monitoramento

### 📈 Status do Ambiente

```bash
# Status geral
./bin/orion-dev status

# Status detalhado
./bin/orion-dev status

# Verificar conectividade
curl http://localhost:7071  # Orion Functions
curl http://localhost:3333  # Orion API
curl http://localhost:10000 # Azurite Storage

# Monitorar recursos
docker stats --no-stream
```

### 📋 Logs

```bash
# Logs de todos os containers
./bin/orion-dev logs

# Logs específicos
./bin/orion-dev logs --service orion-functions
./bin/orion-dev logs --service orion-api
./bin/orion-dev logs --service emulator

# Debug específico
./bin/orion-dev logs --service orion-functions --tail 100
```

### 🌐 URLs de Acesso

| Serviço             | URL                    | Descrição             |
| ------------------- | ---------------------- | --------------------- |
| **Orion Functions** | http://localhost:7071  | Azure Functions local |
| **Orion API**       | http://localhost:3333  | API REST principal    |
| **Azurite Storage** | http://localhost:10000 | Azure Storage local   |
| **Service Bus**     | sb://localhost:5672    | Service Bus Emulator  |
| **PostgreSQL**      | localhost:5432         | Banco de dados        |
| **SQL Server Edge** | localhost:1433         | Banco secundário      |

---

## 🛠️ Desenvolvimento

### 🔄 Workflow de Desenvolvimento

```bash
# 1. Compilar aplicação Go
go build -o bin/orion-dev cmd/main.go

# 2. Iniciar ambiente
./bin/orion-dev start

# 3. Iniciar proxy Service Bus
./bin/orion-dev proxy

# 4. Desenvolver (editar arquivos)
# Os arquivos são sincronizados automaticamente

# 5. Testar mudanças
./bin/orion-dev check-messages

# 6. Enviar mensagens de teste
./bin/orion-dev push-message sbq.pismo.transaction.creation test.json

# 7. Ver logs em tempo real
./bin/orion-dev logs

# 8. Debug específico do Orion Functions
./bin/orion-dev logs --service orion-functions

# 9. Reconstruir containers
./bin/orion-dev build

# 10. Teste rápido
curl http://localhost:7071/cob/test-id
curl http://localhost:3333

# 11. Formatar código Go
go fmt ./...

# 12. Verificar qualidade do código
go vet ./...

# 13. Parar ambiente
./bin/orion-dev stop
```

### ⚙️ Modificar Configurações

```bash
# Editar configuração das Functions
nano local.settings.json

# Editar configuração do Service Bus
nano docker/service-bus/config.json

# Editar docker-compose
nano docker-compose.yml

# Reiniciar após mudanças
./bin/orion-dev stop
./bin/orion-dev start
```

### 🔨 Reconstruir Containers

```bash
# Reconstruir apenas Orion Functions
./bin/orion-dev rebuild-functions

# Reconstruir apenas Orion API
./bin/orion-dev rebuild-api

# Reconstruir tudo
./bin/orion-dev build
```

---

## 🗄️ Banco de Dados

### 🐘 PostgreSQL (Principal)

```bash
# Configuração
Host: orion-database (interno) / localhost (externo)
Porta: 5432
Usuário: <usuario>
Senha: <senha>
Database: <database>
Schema: orionlocal

# Conectar via psql
psql -h localhost -p 5432 -U <usuario> -d <database>

# Conectar via Docker
docker exec -it database psql -U <usuario> -d <database>
```

### 💾 Volumes de Dados

```bash
# Volumes Docker criados
postgres-data      # Dados PostgreSQL
sqledge-data       # Dados SQL Server Edge
azure-storage-data # Dados Azurite Storage

# Localização dos volumes
docker volume ls | grep orion
```

### 📜 Scripts de Inicialização

```bash
# PostgreSQL
docker/database/init-postgres.sql  # Script de inicialização
docker/database/postgres.conf      # Configuração PostgreSQL
docker/database/certs/             # Certificados SSL

# Verificar logs do banco
docker-compose logs -f database
docker-compose logs -f sqledge
```

---

## 🔍 Troubleshooting

### ❌ Problemas Comuns

#### 0. Service Bus não conecta (Proxy não iniciado)

```bash
# Erro: "connection refused" ou "timeout" ou "context deadline exceeded" em comandos de mensagens
# Solução: Iniciar o proxy Service Bus

# 1. Verificar se o proxy está rodando
lsof -i :5671

# 2. Iniciar proxy Service Bus
./bin/orion-dev proxy

# 3. Testar conectividade
./bin/orion-dev check-messages

# 4. Se ainda não funcionar, verificar se o Service Bus Emulator está rodando
docker-compose ps emulator
```

#### 1. Porta já em uso

```bash
# Verificar processos
lsof -i :7071
lsof -i :3333
lsof -i :5672
lsof -i :10000
lsof -i :5432
lsof -i :1433

# Parar processos
sudo kill -9 <PID>
```

#### 2. Containers não iniciam

```bash
# Verificar logs
docker-compose logs

# Verificar status
docker-compose ps

# Reiniciar containers
docker-compose restart
```

#### 3. Functions não conectam

```bash
# Verificar saúde
./bin/orion-dev status

# Verificar logs
docker-compose logs -f orion-functions

# Verificar conectividade
curl http://localhost:7071
```

#### 4. Erro de build

```bash
# Limpar cache
docker system prune -f

# Reconstruir
docker-compose build --no-cache

# Reconstruir específico
docker-compose build orion-functions --no-cache
```

#### 5. Erros de inicialização dos containers

#### 5.1 Erro ao iniciar o container **sqledge**

Caso ocorra erro ao iniciar o container **sqledge**, tente rodar o comando abaixo:

```bash
# Parar e remover apenas o container
docker compose stop sqledge
docker compose rm -f sqledge

# Verificar nome do volume associado
docker volume ls | grep sqledge

# Remover volume explicitamente
docker volume rm finoriondev_sqledge-data

# Reiniciar o container
docker compose up -d sqledge
```

Se necessário reinicie os containers dependentes do **sqledge**, rode o comando abaixo:

```bash
# Parar e iniciar os containers dependentes
docker compose stop orion-functions emulator

# Iniciar os containers dependentes
docker compose start orion-functions emulator
```

#### 6. Problemas com mensagens

```bash
# Verificar arquivos JSON
./bin/orion-dev list

# Validar JSON
./bin/orion-dev validate-json messages/transaction.json

# Formatar JSON
./bin/orion-dev format-json messages/transaction.json

# Mostrar JSON formatado
./bin/orion-dev show-json messages/transaction.json

# Verificar filas
./bin/orion-dev list

# Verificar Service Bus
./bin/orion-dev check-messages
```

#### 7. Problemas com banco de dados

```bash
# Verificar conectividade PostgreSQL
docker exec -it database pg_isready -U postgres

# Verificar logs dos bancos
docker-compose logs database
docker-compose logs sqledge

# Reiniciar apenas os bancos
docker-compose restart database sqledge
```

#### 8. Problemas com Service Bus (Proxy)

```bash
# Verificar se o proxy está rodando
lsof -i :5671

# Verificar se o Service Bus Emulator está rodando
lsof -i :5672

# Iniciar proxy Service Bus
./bin/orion-dev proxy

# Verificar conectividade
telnet localhost 5671
telnet localhost 5672

# Testar comando de mensagens
./bin/orion-dev check-messages

# Se o proxy não iniciar, verificar logs
./bin/orion-dev proxy --verbose
```

### 🧹 Limpeza Completa

```bash
# Parar e remover tudo
./bin/orion-dev stop --clean

# Limpar volumes
docker-compose down -v

# Limpar imagens
docker image prune -f

# Reconstruir
./bin/orion-dev start
```

### 📋 Logs de Debug

```bash
# Debug completo
docker-compose logs -f

# Debug específico
docker-compose logs -f orion-functions

# Logs em tempo real
docker-compose logs -f --tail=100

# Logs de erro
docker-compose logs --tail=50 | grep -i error
```

---

## 🏷️ Releases e Versionamento

### 🚀 Sistema de Releases Automatizado

O projeto utiliza um sistema completo de releases automatizado baseado em **Semantic Versioning (SemVer)**. Quando uma tag é criada, o GitHub Actions automaticamente:

1. ✅ Executa todos os testes
2. 🔨 Compila binários para 6 plataformas
3. 🔐 Gera checksums SHA256
4. 📝 Cria changelog baseado em commits
5. 🚀 Publica o release no GitHub

### 📦 Plataformas Suportadas

Cada release inclui binários para:

| Plataforma              | Arquivo                       | Arquitetura |
| ----------------------- | ----------------------------- | ----------- |
| **macOS Intel**         | `orion-dev-darwin-amd64`      | x86_64      |
| **macOS Apple Silicon** | `orion-dev-darwin-arm64`      | ARM64       |
| **Linux Intel**         | `orion-dev-linux-amd64`       | x86_64      |
| **Linux ARM**           | `orion-dev-linux-arm64`       | ARM64       |
| **Windows Intel**       | `orion-dev-windows-amd64.exe` | x86_64      |
| **Windows ARM**         | `orion-dev-windows-arm64.exe` | ARM64       |

### 🚀 Como Criar um Release

#### Método 1: Script Automatizado (Recomendado)

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

#### Método 2: Via Makefile

```bash
# Criar release
make release VERSION=1.2.3

# Testar release (dry-run)
make release-dry-run VERSION=1.2.3

# Criar release e fazer push
make release-push VERSION=1.2.3

# Compilar para todas as plataformas
make build-all

# Ver versão atual
make version

# Sugerir próxima versão
make version-next
```

#### Método 3: Manual

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

### 🔐 Verificação de Integridade

Cada binário inclui um arquivo `.sha256` para verificação:

```bash
# Download
wget https://github.com/fabriciojbo/fin.orion.dev/releases/download/v1.0.3/orion-dev-linux-amd64
wget https://github.com/fabriciojbo/fin.orion.dev/releases/download/v1.0.3/orion-dev-linux-amd64.sha256

# Verificar
sha256sum -c orion-dev-linux-amd64.sha256

# Tornar executável
chmod +x orion-dev-linux-amd64
```

### 📝 Conventional Commits

Para gerar changelogs automáticos, use commits no formato:

```bash
git commit -m "feat: adicionar comando de proxy TLS"
git commit -m "fix(messages): corrigir validação de JSON"
git commit -m "docs: atualizar README com novos comandos"
```

**Tipos válidos:**
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Documentação
- `style`: Formatação
- `refactor`: Refatoração
- `test`: Testes
- `chore`: Manutenção

### 📊 Monitoramento de Releases

- **Status do Workflow**: Acompanhe em **Actions** no GitHub
- **Tempo de Build**: ~5-10 minutos
- **Cobertura de Testes**: >80%
- **Tamanho dos Binários**: ~10-15MB cada

### 🔄 Atualizações Automáticas

O projeto inclui **Dependabot** configurado para:

- Atualizar dependências Go semanalmente
- Atualizar GitHub Actions semanalmente
- Criar Pull Requests automáticos
- Executar testes antes de mergear

---

## 📚 Documentação

### 📁 Estrutura de Arquivos

- **📁 `cmd/`** - Ponto de entrada da aplicação Go
- **📁 `internal/`** - Código interno da aplicação
  - **📁 `commands/`** - Comandos CLI
  - **📁 `servicebus/`** - Cliente Azure Service Bus
  - **📁 `utils/`** - Utilitários
- **📁 `docker/`** - Configurações Docker e Service Bus
- **📁 `messages/`** - Arquivos JSON de teste
- **📁 `.github/workflows/`** - GitHub Actions
- **📄 `docker-compose.yml`** - Orquestração dos serviços
- **📄 `local.settings.json`** - Configuração Orion Functions
- **📄 `go.mod`** - Dependências Go
- **📄 `go.sum`** - Checksums das dependências
- **📄 `.gitignore`** - Arquivos ignorados pelo Git

### 📋 Comandos CLI Disponíveis

| Comando                          | Função                                    |
| -------------------------------- | ----------------------------------------- |
| `./bin/orion-dev setup`          | Configuração inicial do ambiente          |
| `./bin/orion-dev start`          | Iniciar todos os serviços                 |
| `./bin/orion-dev stop`           | Parar ambiente                            |
| `./bin/orion-dev status`         | Verificar status dos containers           |
| `./bin/orion-dev proxy`          | Iniciar proxy Service Bus (TLS 5671→5672) |
| `./bin/orion-dev list`           | Listar recursos disponíveis               |
| `./bin/orion-dev push-message`   | Enviar mensagem para fila                 |
| `./bin/orion-dev check-messages` | Verificar mensagens do Service Bus        |
| `./bin/orion-dev check-queue`    | Verificar mensagens de fila específica    |
| `./bin/orion-dev check-topic`    | Verificar mensagens do tópico             |
| `./bin/orion-dev validate-json`  | Validar arquivo JSON                      |
| `./bin/orion-dev format-json`    | Formatar arquivo JSON                     |
| `./bin/orion-dev show-json`      | Mostrar JSON formatado                    |

### 📦 Dependências do Projeto

#### Produção
- `github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus` - Cliente Azure Service Bus
- `github.com/spf13/cobra` - Framework CLI
- `github.com/fatih/color` - Cores no terminal
- `github.com/joho/godotenv` - Carregamento de variáveis de ambiente

#### Desenvolvimento
- `go` - Compilador Go 1.24+
- `docker` - Docker Engine
- `docker-compose` - Docker Compose

### 🔗 Recursos Externos

- [Azure Functions Documentation](https://docs.microsoft.com/en-us/azure/azure-functions/)
- [Service Bus Emulator](https://github.com/Azure/azure-sdk-for-net/tree/main/sdk/servicebus/Microsoft.Azure.ServiceBus/Emulator)
- [Azurite Documentation](https://github.com/Azure/Azurite)
- [Azure SQL Edge](https://docs.microsoft.com/en-us/azure/azure-sql-edge/)
- [PostgreSQL Docker](https://hub.docker.com/_/postgres)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go Documentation](https://golang.org/doc/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)

---

## 🎯 Próximos Passos

1. **Compilar aplicação**: `go build -o bin/orion-dev cmd/main.go`
2. **Configurar ambiente**: `./bin/orion-dev setup`
3. **Iniciar serviços**: `./bin/orion-dev start`
4. **Iniciar proxy Service Bus**: `./bin/orion-dev proxy`
5. **Testar funcionalidade**: `./bin/orion-dev check-messages`
6. **Enviar mensagens**: `./bin/orion-dev push-message sbq.pismo.transaction.creation test.json`
7. **Monitorar logs**: `docker-compose logs -f`
8. **Parar ambiente**: `./bin/orion-dev stop`
9. **Criar release**: `./scripts/release.sh -v 1.0.0 -p`

---

**📝 Nota**: Este ambiente é **apenas para desenvolvimento e testes**. Não use em produção. Os dados são voláteis e serão perdidos ao parar os containers.

---

## 📝 Conventional Commits

O projeto utiliza um sistema nativo em Go para validar mensagens de commit seguindo o padrão **Conventional Commits**.

### 🚀 Comandos Rápidos

```bash
# Validar último commit
./bin/orion-dev commitlint-last

# Ver tipos válidos
./bin/orion-dev commitlint-types

# Formatar mensagem
./bin/orion-dev commitlint-format feat auth "adicionar autenticação"

# Instalar hooks automáticos
./scripts/install-hooks.sh
```

### 📋 Exemplos de Commits Válidos

```bash
feat: adicionar comando de proxy TLS
fix(messages): corrigir validação de JSON
docs: atualizar README com novos comandos
test(api): adicionar testes para endpoints
chore(deps): atualizar dependências
feat!: breaking change
fix(auth)!: breaking change com escopo
```

### 🔧 Hooks Automáticos

- **commit-msg**: Valida mensagens automaticamente
- **pre-commit**: Executa testes antes do commit
- **pre-push**: Valida antes do push

Para mais detalhes, consulte [docs/CONVENTIONAL_COMMITS.md](docs/CONVENTIONAL_COMMITS.md).
