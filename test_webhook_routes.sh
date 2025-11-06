#!/bin/bash

# Script de teste para as novas rotas de webhook
# Uso: ./test_webhook_routes.sh

set -e

# Configurações
BASE_URL="http://localhost:8080"
API_KEY="your-api-key-here"
SESSION_ID="test-session-123"

echo "🧪 Testando Rotas de Webhook - zpmeow"
echo "========================================"
echo ""

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${YELLOW}Testando:${NC} $description"
    echo "Método: $method"
    echo "Endpoint: $endpoint"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "apikey: $API_KEY" \
            -H "Content-Type: application/json")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "apikey: $API_KEY" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✅ Sucesso (HTTP $http_code)${NC}"
        echo "Resposta:"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ Erro (HTTP $http_code)${NC}"
        echo "Resposta:"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    fi
    
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Teste 1: Criar sessão de teste (se não existir)
echo "📝 Passo 1: Criar sessão de teste"
test_endpoint "POST" "/sessions/create" \
    '{"name":"Test Session for Webhook"}' \
    "Criar sessão de teste"

# Aguardar um pouco
sleep 1

# Teste 2: Configurar webhook (SET)
echo "📝 Passo 2: Configurar webhook"
test_endpoint "POST" "/sessions/$SESSION_ID/webhook/set" \
    '{
        "enabled": true,
        "url": "https://webhook.site/unique-id",
        "events": ["message", "connected", "disconnected"],
        "token": "Bearer test-token-123"
    }' \
    "Configurar webhook com eventos"

# Teste 3: Consultar webhook (FIND)
echo "📝 Passo 3: Consultar configuração de webhook"
test_endpoint "GET" "/sessions/$SESSION_ID/webhook/find" \
    "" \
    "Obter configuração atual do webhook"

# Teste 4: Atualizar webhook (SET novamente)
echo "📝 Passo 4: Atualizar webhook"
test_endpoint "POST" "/sessions/$SESSION_ID/webhook/set" \
    '{
        "enabled": true,
        "url": "https://webhook.site/another-id",
        "events": ["message", "status", "qr"]
    }' \
    "Atualizar URL e eventos do webhook"

# Teste 5: Desabilitar webhook
echo "📝 Passo 5: Desabilitar webhook"
test_endpoint "POST" "/sessions/$SESSION_ID/webhook/set" \
    '{
        "enabled": false,
        "url": "https://webhook.site/another-id"
    }' \
    "Desabilitar webhook temporariamente"

# Teste 6: Verificar webhook desabilitado
echo "📝 Passo 6: Verificar webhook desabilitado"
test_endpoint "GET" "/sessions/$SESSION_ID/webhook/find" \
    "" \
    "Verificar que webhook está desabilitado"

# Teste 7: Limpar webhook (CLEAR)
echo "📝 Passo 7: Limpar webhook"
test_endpoint "DELETE" "/sessions/$SESSION_ID/webhook/clear" \
    "" \
    "Remover configuração de webhook"

# Teste 8: Verificar webhook limpo
echo "📝 Passo 8: Verificar webhook limpo"
test_endpoint "GET" "/sessions/$SESSION_ID/webhook/find" \
    "" \
    "Verificar que webhook foi removido"

# Teste 9: Testar validação - URL obrigatória
echo "📝 Passo 9: Testar validação (URL obrigatória)"
test_endpoint "POST" "/sessions/$SESSION_ID/webhook/set" \
    '{
        "enabled": true,
        "url": ""
    }' \
    "Tentar habilitar webhook sem URL (deve falhar)"

# Teste 10: Testar validação - Evento inválido
echo "📝 Passo 10: Testar validação (evento inválido)"
test_endpoint "POST" "/sessions/$SESSION_ID/webhook/set" \
    '{
        "enabled": true,
        "url": "https://webhook.site/test",
        "events": ["invalid_event"]
    }' \
    "Tentar configurar com evento inválido (deve falhar)"

echo ""
echo "========================================"
echo -e "${GREEN}✅ Testes concluídos!${NC}"
echo ""
echo "📊 Resumo:"
echo "- 10 testes executados"
echo "- Verifique os resultados acima"
echo ""
echo "💡 Dica: Use 'jq' para formatar JSON (apt-get install jq)"

