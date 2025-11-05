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

# Inicie o PostgreSQL
docker-compose up -d postgres

# Compile o projeto
go build -o bin/zpmeow ./cmd/server/main.go

# Execute
./bin/zpmeow
```

### Verificar Status

```bash
curl http://localhost:8080/health
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

## 🔌 API Endpoints

### Health Check (Sem Autenticação)
```bash
GET /health
```

### Sessões (Requer Autenticação)

- `POST /api/sessions/create` - Criar sessão
- `GET /api/sessions/list` - Listar sessões
- `GET /api/sessions/:id/info` - Detalhes
- `GET /api/sessions/:id/status` - Status detalhado
- `POST /api/sessions/:id/connect` - Conectar
- `POST /api/sessions/:id/disconnect` - Desconectar
- `POST /api/sessions/:id/pair` - Parear com telefone
- `PUT /api/sessions/:id/webhook` - Atualizar webhook
- `DELETE /api/sessions/:id/delete` - Deletar

## 🔐 Autenticação

A API suporta 3 métodos:

1. **Bearer Token**: `Authorization: Bearer your-api-key`
2. **Header**: `X-API-Key: your-api-key`
3. **Query**: `?api_key=your-api-key`

## 📝 Licença

MIT License
