# ========================================
# ZPMeow - WhatsApp Multi-Device API
# Makefile - Comandos Essenciais
# ========================================

# Variáveis
APP_NAME=zpmeow
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(APP_NAME)
MAIN_PATH=./cmd/zpwoot/main.go
DOCS_DIR=docs
GO=go

# Cores
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m

.PHONY: help
help: ## Mostra comandos disponíveis
	@echo "$(GREEN)ZPMeow - WhatsApp Multi-Device API$(NC)"
	@echo "$(YELLOW)Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-18s$(NC) %s\n", $$1, $$2}'

# ========================================
# Build & Run
# ========================================

.PHONY: build
build: ## Compila o projeto
	@echo "$(YELLOW)Compilando...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@$(GO) build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build concluído$(NC)"

.PHONY: run
run: ## Executa o projeto
	@$(GO) run $(MAIN_PATH)

.PHONY: dev
dev: ## Modo desenvolvimento (hot reload)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)Air não instalado. Execute: make install-tools$(NC)"; \
		$(GO) run $(MAIN_PATH); \
	fi

# ========================================
# Dependências
# ========================================

.PHONY: deps
deps: ## Baixa dependências
	@$(GO) mod download

.PHONY: tidy
tidy: ## Organiza dependências
	@$(GO) mod tidy

.PHONY: install-tools
install-tools: ## Instala ferramentas (swag, air)
	@echo "$(YELLOW)Instalando ferramentas...$(NC)"
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@$(GO) install github.com/cosmtrek/air@latest
	@echo "$(GREEN)✅ Ferramentas instaladas$(NC)"

# ========================================
# Documentação
# ========================================

.PHONY: swagger
swagger: ## Gera documentação Swagger
	@echo "$(YELLOW)Gerando Swagger...$(NC)"
	@if command -v swag > /dev/null 2>&1; then \
		swag init -g $(MAIN_PATH) --output $(DOCS_DIR); \
		echo "$(GREEN)✅ Swagger gerado$(NC)"; \
	else \
		echo "$(RED)Swag não instalado. Execute: make install-tools$(NC)"; \
	fi

# ========================================
# Testes
# ========================================

.PHONY: test
test: ## Executa testes
	@$(GO) test -v ./...

.PHONY: test-coverage
test-coverage: ## Testes com cobertura
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Cobertura: coverage.html$(NC)"

# ========================================
# Qualidade de Código
# ========================================

.PHONY: fmt
fmt: ## Formata código
	@$(GO) fmt ./...

.PHONY: vet
vet: ## Analisa código
	@$(GO) vet ./...

# ========================================
# Docker - Produção
# ========================================

.PHONY: docker-up
docker-up: ## Inicia containers (produção)
	@echo "$(YELLOW)Iniciando containers...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)✅ Containers iniciados$(NC)"

.PHONY: docker-down
docker-down: ## Para containers
	@docker compose down

.PHONY: docker-down-v
docker-down-v: ## Para containers e remove volumes
	@echo "$(YELLOW)Parando containers e removendo volumes...$(NC)"
	@docker compose down -v
	@echo "$(GREEN)✅ Containers e volumes removidos$(NC)"

.PHONY: docker-logs
docker-logs: ## Logs dos containers
	@docker compose logs -f

.PHONY: docker-build
docker-build: ## Build imagem Docker
	@docker compose build

# ========================================
# Docker - Desenvolvimento
# ========================================

.PHONY: dev-up
dev-up: ## Inicia ambiente dev (postgres+nats+dbgate+webhook)
	@echo "$(YELLOW)Iniciando ambiente de desenvolvimento...$(NC)"
	@docker compose -f docker-compose.dev.yml up -d
	@echo "$(GREEN)✅ Ambiente dev iniciado$(NC)"
	@echo "$(YELLOW)Serviços disponíveis:$(NC)"
	@echo "  PostgreSQL:      localhost:5432"
	@echo "  NATS:            localhost:4222"
	@echo "  NATS Monitor:    http://localhost:8222"
	@echo "  DBGate:          http://localhost:3000"
	@echo "  Webhook Tester:  http://localhost:8090"

.PHONY: dev-down
dev-down: ## Para ambiente dev
	@docker compose -f docker-compose.dev.yml down

.PHONY: dev-down-v
dev-down-v: ## Para ambiente dev e remove volumes
	@echo "$(YELLOW)Parando ambiente dev e removendo volumes...$(NC)"
	@docker compose -f docker-compose.dev.yml down -v
	@echo "$(GREEN)✅ Ambiente dev e volumes removidos$(NC)"

.PHONY: dev-logs
dev-logs: ## Logs do ambiente dev
	@docker compose -f docker-compose.dev.yml logs -f

.PHONY: dev-restart
dev-restart: dev-down dev-up ## Reinicia ambiente dev

# ========================================
# Serviços Individuais
# ========================================

.PHONY: db-up
db-up: ## Inicia PostgreSQL
	@docker compose -f docker-compose.dev.yml up -d postgres
	@echo "$(GREEN)✅ PostgreSQL: localhost:5432$(NC)"

.PHONY: nats-up
nats-up: ## Inicia NATS
	@docker compose -f docker-compose.dev.yml up -d nats
	@echo "$(GREEN)✅ NATS: localhost:4222$(NC)"

.PHONY: dbgate-up
dbgate-up: ## Inicia DBGate
	@docker compose -f docker-compose.dev.yml up -d dbgate
	@echo "$(GREEN)✅ DBGate: http://localhost:3000$(NC)"

.PHONY: webhook-up
webhook-up: ## Inicia Webhook Tester
	@docker compose -f docker-compose.dev.yml up -d webhook-tester
	@echo "$(GREEN)✅ Webhook Tester: http://localhost:8090$(NC)"

# ========================================
# Utilitários
# ========================================

.PHONY: clean
clean: ## Remove arquivos compilados
	@rm -rf $(BINARY_DIR) coverage.out coverage.html

.PHONY: clean-all
clean-all: ## Remove TUDO (containers, volumes, binários)
	@echo "$(RED)⚠️  Removendo TODOS os containers e volumes...$(NC)"
	@docker compose down -v 2>/dev/null || true
	@docker compose -f docker-compose.dev.yml down -v 2>/dev/null || true
	@rm -rf $(BINARY_DIR) coverage.out coverage.html
	@echo "$(GREEN)✅ Limpeza completa concluída$(NC)"

.PHONY: kill
kill: ## Mata processo na porta 8080
	@if command -v lsof > /dev/null 2>&1; then \
		lsof -ti:8080 2>/dev/null | xargs -r kill -9 && echo "$(GREEN)✅ Processo finalizado$(NC)" || echo "$(GREEN)✅ Porta livre$(NC)"; \
	else \
		fuser -k 8080/tcp 2>/dev/null && echo "$(GREEN)✅ Processo finalizado$(NC)" || echo "$(GREEN)✅ Porta livre$(NC)"; \
	fi

.PHONY: health
health: ## Verifica saúde da API
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)API não está respondendo$(NC)"

# ========================================
# Comandos Compostos
# ========================================

.PHONY: setup
setup: deps install-tools swagger ## Setup inicial completo
	@echo "$(GREEN)✅ Setup completo!$(NC)"
	@echo "$(YELLOW)Próximos passos:$(NC)"
	@echo "  1. make dev-up    # Inicia ambiente dev"
	@echo "  2. make run       # Executa a API"

.PHONY: all
all: clean tidy swagger build ## Build completo

# Default
.DEFAULT_GOAL := help

