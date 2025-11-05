# zpmeow - WhatsApp REST API

API REST para integração com WhatsApp usando a biblioteca whatsmeow.

## Visão Geral

Este projeto expõe funcionalidades do WhatsApp via API REST. Principais recursos planejados: envio de mensagens, gerenciamento de grupos, contatos e sessões.

Tecnologias utilizadas: Go, Gin, GORM, whatsmeow.

## Estrutura do Projeto

```
zpmeow/
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

## Requisitos

- Go 1.21 ou superior
- SQLite (desenvolvimento) ou PostgreSQL (produção)
- Conta WhatsApp válida para autenticação

## Instalação

```bash
# Clonar o repositório
git clone <repo-url>
cd zpmeow

# Copiar arquivo de configuração
cp .env.example .env

# Editar .env com suas configurações

# Baixar dependências
go mod download
```

## Configuração

Veja `.env.example` para as variáveis principais: PORT, DATABASE_URL, API_KEY, WHATSAPP_DATA_DIR.

## Execução

```bash
# Modo desenvolvimento
go run cmd/server/main.go

# Build para produção
go build -o server cmd/server/main.go
./server
```

## Endpoints da API

A documentação completa será adicionada após a implementação dos handlers. Categorias principais: Sessões, Mensagens, Grupos, Contatos.

## Desenvolvimento

Status: Em desenvolvimento

Próximas fases: Implementação de handlers, services e repository.

## Licença

Escolha uma licença (por exemplo MIT) e adicione-a neste arquivo.
