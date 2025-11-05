# 📚 Guia de Uso do Swagger - ZPWoot API

## 🚀 Acesso Rápido

Após iniciar o servidor, acesse a documentação interativa em:

```
http://localhost:8080/swagger/index.html
```

## 🔐 Autenticação no Swagger

A API suporta 3 métodos de autenticação. No Swagger UI, você pode usar qualquer um deles:

### Método 1: Bearer Token (Recomendado)
1. Clique no botão **"Authorize"** no topo da página
2. No campo **BearerAuth**, digite: `Bearer SEU_TOKEN_AQUI`
3. Clique em **"Authorize"** e depois **"Close"**

### Método 2: X-API-Key Header
1. Clique no botão **"Authorize"**
2. No campo **ApiKeyAuth**, digite apenas: `SEU_TOKEN_AQUI`
3. Clique em **"Authorize"** e depois **"Close"**

### Método 3: Query Parameter
1. Clique no botão **"Authorize"**
2. No campo **ApiKeyQuery**, digite: `SEU_TOKEN_AQUI`
3. Clique em **"Authorize"** e depois **"Close"**

> **Nota**: O token é configurado na variável `API_KEY` no arquivo `.env`

## 📝 Testando Endpoints

### 1. Health Check (Sem autenticação)
- Endpoint: `GET /health`
- Clique em **"Try it out"**
- Clique em **"Execute"**
- Você deve ver: `{"status": "ok", "service": "zpwoot"}`

### 2. Criar uma Sessão
- Endpoint: `POST /sessions/create`
- Clique em **"Try it out"**
- Edite o JSON de exemplo:
```json
{
  "name": "Minha Primeira Sessão",
  "webhook_url": "https://seu-webhook.com/whatsapp",
  "webhook_events": ["message", "qr"],
  "metadata": {
    "cliente": "Empresa XYZ"
  }
}
```
- Clique em **"Execute"**
- Copie o `id` da sessão retornada

### 3. Listar Sessões
- Endpoint: `GET /sessions/list`
- Clique em **"Try it out"**
- Clique em **"Execute"**
- Veja todas as sessões criadas

### 4. Parear com Telefone
- Endpoint: `POST /sessions/{id}/pair`
- Clique em **"Try it out"**
- Cole o `id` da sessão no campo **id**
- Edite o JSON:
```json
{
  "phone_number": "+5511999999999"
}
```
- Clique em **"Execute"**
- Você receberá um código de 8 dígitos
- Digite esse código no WhatsApp do seu celular

### 5. Verificar Status da Sessão
- Endpoint: `GET /sessions/{id}/status`
- Clique em **"Try it out"**
- Cole o `id` da sessão
- Clique em **"Execute"**
- Veja o status detalhado da conexão

## 🎯 Recursos do Swagger UI

### Modelos (Schemas)
- Role até o final da página para ver todos os modelos de dados
- Clique em cada modelo para expandir e ver todos os campos
- Útil para entender a estrutura de requisições e respostas

### Filtros por Tag
- Use as tags (ex: "Sessions") para filtrar endpoints
- Clique na tag para expandir/colapsar todos os endpoints daquela categoria

### Download da Especificação
- Acesse `http://localhost:8080/swagger/doc.json` para JSON
- Acesse `http://localhost:8080/swagger/doc.yaml` para YAML
- Use para importar em outras ferramentas (Postman, Insomnia, etc.)

## 🔄 Regenerar Documentação

Se você modificar os handlers ou adicionar novos endpoints:

```bash
# Certifique-se de ter o swag instalado
go install github.com/swaggo/swag/cmd/swag@latest

# Regenere a documentação
swag init -g cmd/zpwoot/main.go --output docs

# Recompile e reinicie o servidor
go build -o bin/zpmeow ./cmd/zpwoot/main.go
./bin/zpmeow
```

## 📖 Anotações Swagger nos Handlers

Exemplo de como documentar um novo endpoint:

```go
// MinhaFuncao faz algo incrível
// @Summary Resumo curto
// @Description Descrição detalhada do que o endpoint faz
// @Tags NomeDaTag
// @Accept json
// @Produce json
// @Param id path string true "ID do recurso"
// @Param request body dto.MeuRequest true "Dados da requisição"
// @Success 200 {object} dto.MeuResponse
// @Failure 400 {object} dto.ErrorResponse
// @Security BearerAuth
// @Security ApiKeyAuth
// @Security ApiKeyQuery
// @Router /meu-endpoint/{id} [post]
func (h *Handler) MinhaFuncao(c *gin.Context) {
    // implementação
}
```

## 🆘 Problemas Comuns

### Erro 401 Unauthorized
- Verifique se você autenticou usando o botão "Authorize"
- Confirme que o token está correto no arquivo `.env`
- Certifique-se de incluir "Bearer " antes do token (se usar BearerAuth)

### Swagger não carrega
- Verifique se o servidor está rodando: `curl http://localhost:8080/health`
- Confirme que a porta 8080 está livre
- Verifique os logs do servidor para erros

### Documentação desatualizada
- Execute `swag init` novamente
- Recompile o projeto
- Reinicie o servidor

