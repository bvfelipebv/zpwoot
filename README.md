# 🐱 zpmeow - WhatsApp Multi-Device API

API REST completa para gerenciar múltiplas sessões WhatsApp usando whatsmeow.

## ✨ Características

- 🔄 **Múltiplas Sessões**: Gerencie várias contas WhatsApp simultaneamente
- 🔐 **Autenticação Segura**: API Key com 3 métodos de autenticação
- 📱 **Pareamento Flexível**: QR Code ou código via telefone
- 🗄️ **PostgreSQL**: Banco de dados robusto e escalável
- 🔄 **Auto-Restore**: Reconecta sessões automaticamente ao reiniciar
- 📡 **Webhooks**: Receba eventos em tempo real
- 🚀 **Graceful Shutdown**: Desconexão segura de todas as sessões
- 📝 **Logs Estruturados**: Zerolog para logging profissional

## Estrutura do Projeto

```
zpwoot/
├── cmd/server/          # Ponto de entrada da aplicação
├── internal/
│   ├── api/            # Handlers HTTP e routers
│   ├── config/         # Configurações da aplicação
│   ├── model/          # Estruturas de dados
│   ├── repository/     # Camada de acesso a dados
│   └── service/        # Lógica de negócio
├── pkg/
│   ├── logger/         # Sistema de logging
│   └── utils/          # Funções utilitárias
└── go.mod              # Dependências do projeto
```

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 16 (via Docker)

### Instalação

```bash
# Clone o repositório
git clone https://github.com/bvfelipebv/zpwoot.git
cd zpwoot

# Copie o arquivo de configuração
cp .env.example .env

# Edite o .env e configure sua API_KEY
nano .env

# Setup inicial (instala ferramentas e gera documentação)
make setup

# Inicie o PostgreSQL
make db-up

# Compile e execute
make start
```

### Comandos Make Disponíveis

```bash
# Ver todos os comandos disponíveis
make help

# Comandos principais
make build          # Compila o projeto
make run            # Executa sem compilar
make start          # Compila e executa
make kill           # Mata processos na porta 8080
make swagger        # Gera documentação Swagger

# Docker
make db-up          # Inicia PostgreSQL
make db-down        # Para PostgreSQL
make docker-up      # Inicia todos os containers
make docker-down    # Para todos os containers

# Desenvolvimento
make dev            # Modo desenvolvimento com hot reload
make test           # Executa testes
make fmt            # Formata código
make clean          # Limpa arquivos compilados
```

### Verificar Status

```bash
curl http://localhost:8080/health
# ou
make health
```

## 🐳 Docker

```bash
# Iniciar PostgreSQL
docker-compose up -d postgres

# Iniciar DBGate (interface web para gerenciar o banco)
docker-compose up -d dbgate
# Acesse: http://localhost:3000

# Iniciar todos os serviços
docker-compose up -d
```

### Interface de Gerenciamento

**DBGate** - http://localhost:3000
- ✅ Interface moderna e intuitiva
- ✅ Conexão pré-configurada automaticamente
- ✅ Query builder visual
- ✅ Importação/exportação de dados
- ✅ Suporte a múltiplos bancos de dados
- ✅ Sem necessidade de configuração manual

## 📚 Documentação da API (Swagger)

A documentação interativa da API está disponível via Swagger UI:

```
http://localhost:8080/swagger/index.html
```

### Recursos do Swagger:
- ✅ Documentação completa de todos os endpoints
- ✅ Teste interativo de APIs diretamente no navegador
- ✅ **Exemplos completos** de requisições e respostas
- ✅ Modelos de dados detalhados com valores de exemplo
- ✅ Autenticação integrada (API Key)
- ✅ Host dinâmico (funciona em qualquer ambiente)

### Regenerar Documentação Swagger

Se você fizer alterações nos handlers ou adicionar novos endpoints:

```bash
# Instalar swag CLI (se ainda não tiver)
go install github.com/swaggo/swag/cmd/swag@latest

# Regenerar documentação
swag init -g cmd/zpwoot/main.go --output docs
```

## 🔌 API Endpoints

### Health Check (Sem Autenticação)
```bash
GET /health
```

### Sessões (Requer Autenticação)

- `POST /sessions/create` - Criar sessão
- `GET /sessions/list` - Listar sessões
- `GET /sessions/:id/info` - Detalhes
- `GET /sessions/:id/status` - Status detalhado
- `POST /sessions/:id/connect` - Conectar
- `POST /sessions/:id/disconnect` - Desconectar
- `POST /sessions/:id/pair` - Parear com telefone
- `PUT /sessions/:id/webhook` - Atualizar webhook
- `DELETE /sessions/:id/delete` - Deletar

## 🔐 Autenticação

A API usa autenticação simples via header:

**Header**: `apikey: your-api-key`

**Exemplo:**
```bash
curl -H "apikey: sldkfjsldkflskdfjlsd" http://localhost:8080/sessions/list
```

Configure sua API Key no arquivo `.env`:
```bash
API_KEY=sldkfjsldkflskdfjlsd
```

## 📝 Licença

MIT License
