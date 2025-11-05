# ========================================
# ZPWoot - WhatsApp Multi-Device API
# Makefile para comandos básicos
# ========================================

# Variáveis
APP_NAME=zpmeow
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(APP_NAME)
MAIN_PATH=./cmd/zpwoot/main.go
DOCS_DIR=docs
GO=go
SWAG=swag

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help
help: ## Mostra esta mensagem de ajuda
	@echo "$(GREEN)ZPWoot - WhatsApp Multi-Device API$(NC)"
	@echo "$(YELLOW)Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# ========================================
# Comandos de Build
# ========================================

.PHONY: build
build: ## Compila o projeto
	@echo "$(YELLOW)Compilando $(APP_NAME)...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@$(GO) build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build concluído: $(BINARY_PATH)$(NC)"

.PHONY: build-linux
build-linux: ## Compila para Linux
	@echo "$(YELLOW)Compilando para Linux...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Build Linux concluído$(NC)"

.PHONY: build-windows
build-windows: ## Compila para Windows
	@echo "$(YELLOW)Compilando para Windows...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)✅ Build Windows concluído$(NC)"

.PHONY: build-mac
build-mac: ## Compila para macOS
	@echo "$(YELLOW)Compilando para macOS...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "$(GREEN)✅ Build macOS concluído$(NC)"

.PHONY: build-all
build-all: build-linux build-windows build-mac ## Compila para todas as plataformas
	@echo "$(GREEN)✅ Build completo para todas as plataformas$(NC)"

# ========================================
# Comandos de Execução
# ========================================

.PHONY: run
run: ## Executa o projeto
	@echo "$(YELLOW)Executando $(APP_NAME)...$(NC)"
	@$(GO) run $(MAIN_PATH)

.PHONY: start
start: build ## Compila e executa o projeto
	@echo "$(YELLOW)Iniciando $(APP_NAME)...$(NC)"
	@./$(BINARY_PATH)

# ========================================
# Comandos de Desenvolvimento
# ========================================

.PHONY: dev
dev: ## Executa em modo desenvolvimento com hot reload (requer air)
	@echo "$(YELLOW)Iniciando em modo desenvolvimento...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)❌ Air não instalado. Instale com: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(YELLOW)Executando sem hot reload...$(NC)"; \
		$(GO) run $(MAIN_PATH); \
	fi

.PHONY: install-tools
install-tools: ## Instala ferramentas de desenvolvimento
	@echo "$(YELLOW)Instalando ferramentas...$(NC)"
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@$(GO) install github.com/cosmtrek/air@latest
	@echo "$(GREEN)✅ Ferramentas instaladas$(NC)"

# ========================================
# Comandos de Dependências
# ========================================

.PHONY: deps
deps: ## Baixa as dependências do projeto
	@echo "$(YELLOW)Baixando dependências...$(NC)"
	@$(GO) mod download
	@echo "$(GREEN)✅ Dependências baixadas$(NC)"

.PHONY: tidy
tidy: ## Organiza as dependências
	@echo "$(YELLOW)Organizando dependências...$(NC)"
	@$(GO) mod tidy
	@echo "$(GREEN)✅ Dependências organizadas$(NC)"

.PHONY: vendor
vendor: ## Cria pasta vendor com dependências
	@echo "$(YELLOW)Criando vendor...$(NC)"
	@$(GO) mod vendor
	@echo "$(GREEN)✅ Vendor criado$(NC)"

# ========================================
# Comandos de Documentação
# ========================================

.PHONY: swagger
swagger: ## Gera documentação Swagger
	@echo "$(YELLOW)Gerando documentação Swagger...$(NC)"
	@if command -v swag > /dev/null 2>&1; then \
		swag init -g $(MAIN_PATH) --output $(DOCS_DIR); \
		echo "$(GREEN)✅ Documentação Swagger gerada em $(DOCS_DIR)$(NC)"; \
	elif [ -f $(HOME)/go/bin/swag ]; then \
		$(HOME)/go/bin/swag init -g $(MAIN_PATH) --output $(DOCS_DIR); \
		echo "$(GREEN)✅ Documentação Swagger gerada em $(DOCS_DIR)$(NC)"; \
	else \
		echo "$(RED)❌ Swag não instalado. Execute: make install-tools$(NC)"; \
	fi

.PHONY: swagger-fmt
swagger-fmt: ## Formata anotações Swagger
	@echo "$(YELLOW)Formatando anotações Swagger...$(NC)"
	@swag fmt
	@echo "$(GREEN)✅ Anotações formatadas$(NC)"

# ========================================
# Comandos de Testes
# ========================================

.PHONY: test
test: ## Executa os testes
	@echo "$(YELLOW)Executando testes...$(NC)"
	@$(GO) test -v ./...

.PHONY: test-coverage
test-coverage: ## Executa testes com cobertura
	@echo "$(YELLOW)Executando testes com cobertura...$(NC)"
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Relatório de cobertura gerado: coverage.html$(NC)"

.PHONY: test-race
test-race: ## Executa testes com detector de race condition
	@echo "$(YELLOW)Executando testes com race detector...$(NC)"
	@$(GO) test -race -v ./...

# ========================================
# Comandos de Qualidade de Código
# ========================================

.PHONY: fmt
fmt: ## Formata o código
	@echo "$(YELLOW)Formatando código...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✅ Código formatado$(NC)"

.PHONY: vet
vet: ## Analisa o código com go vet
	@echo "$(YELLOW)Analisando código...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✅ Análise concluída$(NC)"

.PHONY: lint
lint: ## Executa linter (requer golangci-lint)
	@echo "$(YELLOW)Executando linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✅ Linting concluído$(NC)"; \
	else \
		echo "$(RED)❌ golangci-lint não instalado$(NC)"; \
		echo "$(YELLOW)Instale com: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin$(NC)"; \
	fi

.PHONY: check
check: fmt vet ## Formata e analisa o código
	@echo "$(GREEN)✅ Verificações concluídas$(NC)"

# ========================================
# Comandos de Docker
# ========================================

.PHONY: docker-build
docker-build: ## Constrói a imagem Docker
	@echo "$(YELLOW)Construindo imagem Docker...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)✅ Imagem Docker construída$(NC)"

.PHONY: docker-up
docker-up: ## Inicia os containers Docker
	@echo "$(YELLOW)Iniciando containers...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)✅ Containers iniciados$(NC)"

.PHONY: docker-down
docker-down: ## Para os containers Docker
	@echo "$(YELLOW)Parando containers...$(NC)"
	@docker compose down
	@echo "$(GREEN)✅ Containers parados$(NC)"

.PHONY: docker-logs
docker-logs: ## Mostra logs dos containers
	@docker compose logs -f

.PHONY: docker-ps
docker-ps: ## Lista containers em execução
	@docker compose ps

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Reinicia os containers

.PHONY: db-up
db-up: ## Inicia apenas o PostgreSQL
	@echo "$(YELLOW)Iniciando PostgreSQL...$(NC)"
	@docker compose up -d postgres
	@echo "$(GREEN)✅ PostgreSQL iniciado$(NC)"

.PHONY: db-down
db-down: ## Para o PostgreSQL
	@echo "$(YELLOW)Parando PostgreSQL...$(NC)"
	@docker compose stop postgres
	@echo "$(GREEN)✅ PostgreSQL parado$(NC)"

.PHONY: db-logs
db-logs: ## Mostra logs do PostgreSQL
	@docker compose logs -f postgres

.PHONY: dbgate-up
dbgate-up: ## Inicia DBGate (interface web do banco)
	@echo "$(YELLOW)Iniciando DBGate...$(NC)"
	@docker compose up -d dbgate
	@echo "$(GREEN)✅ DBGate disponível em http://localhost:3000$(NC)"

# ========================================
# Comandos de Limpeza
# ========================================

.PHONY: clean
clean: ## Remove arquivos compilados
	@echo "$(YELLOW)Limpando arquivos...$(NC)"
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✅ Limpeza concluída$(NC)"

.PHONY: clean-all
clean-all: clean ## Remove arquivos compilados e cache
	@echo "$(YELLOW)Limpeza completa...$(NC)"
	@$(GO) clean -cache -testcache -modcache
	@rm -rf vendor
	@echo "$(GREEN)✅ Limpeza completa concluída$(NC)"

# ========================================
# Comandos Úteis
# ========================================

.PHONY: version
version: ## Mostra versão do Go
	@$(GO) version

.PHONY: info
info: ## Mostra informações do projeto
	@echo "$(GREEN)ZPWoot - WhatsApp Multi-Device API$(NC)"
	@echo "$(YELLOW)Informações do Projeto:$(NC)"
	@echo "  App Name:     $(APP_NAME)"
	@echo "  Binary Path:  $(BINARY_PATH)"
	@echo "  Main Path:    $(MAIN_PATH)"
	@echo "  Docs Dir:     $(DOCS_DIR)"
	@echo ""
	@echo "$(YELLOW)Versão do Go:$(NC)"
	@$(GO) version
	@echo ""
	@echo "$(YELLOW)Módulos:$(NC)"
	@$(GO) list -m all | head -5

.PHONY: health
health: ## Verifica se a API está rodando
	@echo "$(YELLOW)Verificando saúde da API...$(NC)"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)❌ API não está respondendo$(NC)"

.PHONY: swagger-ui
swagger-ui: ## Abre o Swagger UI no navegador
	@echo "$(YELLOW)Abrindo Swagger UI...$(NC)"
	@if command -v xdg-open > /dev/null; then \
		xdg-open http://localhost:8080/swagger/index.html; \
	elif command -v open > /dev/null; then \
		open http://localhost:8080/swagger/index.html; \
	else \
		echo "$(YELLOW)Acesse: http://localhost:8080/swagger/index.html$(NC)"; \
	fi

.PHONY: kill
kill: ## Mata processos rodando na porta 8080
	@echo "$(YELLOW)Procurando processos na porta 8080...$(NC)"
	@if command -v lsof > /dev/null 2>&1; then \
		PID=$$(lsof -ti:8080 2>/dev/null); \
		if [ -z "$$PID" ]; then \
			echo "$(GREEN)✅ Nenhum processo rodando na porta 8080$(NC)"; \
		else \
			echo "$(YELLOW)Matando processo(s): $$PID$(NC)"; \
			kill -9 $$PID && echo "$(GREEN)✅ Processo(s) finalizado(s)$(NC)" || echo "$(RED)❌ Erro ao matar processo$(NC)"; \
		fi; \
	elif command -v fuser > /dev/null 2>&1; then \
		if fuser 8080/tcp > /dev/null 2>&1; then \
			echo "$(YELLOW)Matando processo na porta 8080...$(NC)"; \
			fuser -k 8080/tcp && echo "$(GREEN)✅ Processo finalizado$(NC)" || echo "$(RED)❌ Erro ao matar processo$(NC)"; \
		else \
			echo "$(GREEN)✅ Nenhum processo rodando na porta 8080$(NC)"; \
		fi; \
	else \
		PID=$$(ss -lptn 'sport = :8080' 2>/dev/null | grep -oP 'pid=\K[0-9]+' | head -1); \
		if [ -z "$$PID" ]; then \
			PID=$$(netstat -tlnp 2>/dev/null | grep ':8080' | awk '{print $$7}' | cut -d'/' -f1 | head -1); \
		fi; \
		if [ -z "$$PID" ]; then \
			echo "$(GREEN)✅ Nenhum processo rodando na porta 8080$(NC)"; \
		else \
			echo "$(YELLOW)Matando processo: $$PID$(NC)"; \
			kill -9 $$PID && echo "$(GREEN)✅ Processo finalizado$(NC)" || echo "$(RED)❌ Erro ao matar processo$(NC)"; \
		fi; \
	fi

.PHONY: port-check
port-check: ## Verifica qual processo está usando a porta 8080
	@echo "$(YELLOW)Verificando porta 8080...$(NC)"
	@if command -v lsof > /dev/null 2>&1; then \
		lsof -i:8080 2>/dev/null || echo "$(GREEN)✅ Porta 8080 está livre$(NC)"; \
	elif command -v ss > /dev/null 2>&1; then \
		ss -lptn 'sport = :8080' 2>/dev/null || echo "$(GREEN)✅ Porta 8080 está livre$(NC)"; \
	else \
		netstat -tlnp 2>/dev/null | grep ':8080' || echo "$(GREEN)✅ Porta 8080 está livre$(NC)"; \
	fi

# ========================================
# Comandos Compostos
# ========================================

.PHONY: setup
setup: deps install-tools swagger ## Configuração inicial do projeto
	@echo "$(GREEN)✅ Setup completo!$(NC)"
	@echo "$(YELLOW)Próximos passos:$(NC)"
	@echo "  1. Configure o arquivo .env"
	@echo "  2. Execute: make db-up"
	@echo "  3. Execute: make run"

.PHONY: all
all: clean deps swagger build ## Executa limpeza, deps, swagger e build
	@echo "$(GREEN)✅ Build completo concluído!$(NC)"

.PHONY: rebuild
rebuild: clean build ## Limpa e reconstrói o projeto
	@echo "$(GREEN)✅ Rebuild concluído!$(NC)"

.PHONY: deploy-prep
deploy-prep: clean deps test swagger build-all ## Prepara para deploy
	@echo "$(GREEN)✅ Preparação para deploy concluída!$(NC)"

# Default target
.DEFAULT_GOAL := help

