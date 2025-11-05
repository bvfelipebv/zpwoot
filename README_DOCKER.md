# 🐳 Docker Setup - zpmeow

## Visão Geral

Este projeto usa **PostgreSQL** como banco de dados principal e **database/sql** nativo (mesma tecnologia do whatsmeow) ao invés de ORMs.

## 🚀 Início Rápido

### 1. Iniciar o PostgreSQL

```bash
# Iniciar apenas o PostgreSQL
docker-compose up -d postgres

# Verificar se está rodando
docker-compose ps
```

### 2. Configurar variáveis de ambiente

```bash
# Copiar o arquivo de exemplo
cp .env.example .env

# Editar o .env e ajustar as configurações
nano .env
```

### 3. Executar a aplicação localmente

```bash
# A aplicação se conectará ao PostgreSQL no Docker
go run cmd/server/main.go
```

## 📊 Serviços Disponíveis

### PostgreSQL
- **Porta**: 5432
- **Database**: zpmeow
- **User**: zpmeow
- **Password**: zpmeow_password_change_in_production
- **Connection String**: `postgres://zpmeow:zpmeow_password_change_in_production@localhost:5432/zpmeow?sslmode=disable`

### pgAdmin (Interface Web)
- **URL**: http://localhost:5050
- **Email**: admin@zpmeow.local
- **Password**: admin

Para conectar ao PostgreSQL no pgAdmin:
1. Acesse http://localhost:5050
2. Adicione um novo servidor
3. **Host**: postgres (nome do container)
4. **Port**: 5432
5. **Database**: zpmeow
6. **Username**: zpmeow
7. **Password**: zpmeow_password_change_in_production

## 🗄️ Arquitetura de Banco de Dados

### Abordagem: database/sql Nativo

Seguimos a mesma abordagem do whatsmeow:
- ✅ Usa `database/sql` diretamente
- ✅ Compartilha a mesma conexão SQL com o whatsmeow
- ✅ Usa `sqlstore.Container` do whatsmeow para gerenciar devices
- ✅ Driver PostgreSQL: `github.com/lib/pq`
- ❌ **NÃO usa GORM** ou outros ORMs

### Estrutura de Tabelas

#### Nossa Tabela: `sessions`
```sql
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    jid VARCHAR(64) UNIQUE,
    status VARCHAR(20) DEFAULT 'disconnected',
    pairing_method VARCHAR(10),
    push_name VARCHAR(100),
    platform VARCHAR(50),
    business_name VARCHAR(100),
    webhook_url VARCHAR(500),
    webhook_events TEXT,
    webhook_secret VARCHAR(100),
    last_connected TIMESTAMP,
    last_disconnected TIMESTAMP,
    last_paired TIMESTAMP,
    qr_code TEXT,
    qr_code_expiry TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

#### Tabelas do Whatsmeow (criadas automaticamente)
- `whatsmeow_device` - Informações do dispositivo
- `whatsmeow_identity_keys` - Chaves de identidade
- `whatsmeow_pre_keys` - Pre-keys para criptografia
- `whatsmeow_sessions` - Sessões criptográficas
- `whatsmeow_sender_keys` - Chaves de grupo
- `whatsmeow_app_state_sync_keys` - Chaves de sincronização
- `whatsmeow_contacts` - Contatos
- `whatsmeow_chat_settings` - Configurações de chat
- `whatsmeow_message_secrets` - Segredos de mensagens

## 🔧 Comandos Úteis

### Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Parar todos os serviços
docker-compose down

# Ver logs
docker-compose logs -f postgres

# Reiniciar PostgreSQL
docker-compose restart postgres

# Remover volumes (CUIDADO: apaga dados!)
docker-compose down -v
```

### PostgreSQL

```bash
# Conectar ao PostgreSQL via psql
docker-compose exec postgres psql -U zpmeow -d zpmeow

# Backup do banco
docker-compose exec postgres pg_dump -U zpmeow zpmeow > backup.sql

# Restaurar backup
docker-compose exec -T postgres psql -U zpmeow zpmeow < backup.sql
```

## 🔐 Segurança em Produção

**IMPORTANTE**: Antes de ir para produção:

1. **Altere as senhas**:
   - PostgreSQL password
   - pgAdmin password
   - API_KEY no .env

2. **Use SSL/TLS**:
   - Configure `sslmode=require` na connection string
   - Configure certificados SSL no PostgreSQL

3. **Firewall**:
   - Não exponha a porta 5432 publicamente
   - Use VPN ou SSH tunneling

4. **Backups**:
   - Configure backups automáticos
   - Teste restauração regularmente

## 📝 Notas Importantes

1. **Mesma Conexão SQL**: O whatsmeow e nossa aplicação compartilham a mesma conexão `*sql.DB`
2. **Migrations**: O whatsmeow cria suas tabelas automaticamente via `container.Upgrade()`
3. **Transactions**: Use `db.BeginTx()` para transações que envolvem ambas as tabelas
4. **Context**: Sempre passe `context.Context` para operações de banco de dados

## 🐛 Troubleshooting

### Erro: "connection refused"
```bash
# Verificar se o PostgreSQL está rodando
docker-compose ps

# Ver logs do PostgreSQL
docker-compose logs postgres
```

### Erro: "database does not exist"
```bash
# Recriar o banco
docker-compose down
docker-compose up -d postgres
```

### Erro: "too many connections"
```bash
# Aumentar max_connections no PostgreSQL
# Editar docker-compose.yml e adicionar:
# command: postgres -c max_connections=200
```

