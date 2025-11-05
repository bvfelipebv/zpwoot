# Database Migrations

## Visão Geral

Este diretório contém as migrações de banco de dados do zpwoot. As migrações são executadas **automaticamente** na inicialização do sistema.

## Como Funciona

### Execução Automática

Quando você inicia o servidor com `go run cmd/server/main.go`, o sistema:

1. ✅ Conecta ao banco de dados PostgreSQL
2. ✅ Cria a tabela `schema_migrations` (se não existir)
3. ✅ Verifica quais migrações já foram aplicadas
4. ✅ Executa apenas as migrações pendentes em ordem
5. ✅ Registra cada migração aplicada com timestamp

### Estrutura dos Arquivos

Cada migração consiste em dois arquivos:

```
XXX_nome_da_migracao.up.sql    # SQL para aplicar a migração
XXX_nome_da_migracao.down.sql  # SQL para reverter a migração
```

Onde `XXX` é o número da versão (ex: 001, 002, 003...)

## Migrações Existentes

### 001_create_sessions

Cria a tabela principal `sessions` com:
- ID como UUID (gerado automaticamente)
- Campos para gerenciar sessões WhatsApp
- Índices para performance
- Trigger para atualizar `updated_at` automaticamente

## Como Criar uma Nova Migração

### 1. Criar os arquivos

```bash
# Exemplo: Adicionar campo "last_seen" na tabela sessions
touch internal/db/migrations/002_add_last_seen.up.sql
touch internal/db/migrations/002_add_last_seen.down.sql
```

### 2. Escrever o SQL de UP (002_add_last_seen.up.sql)

```sql
-- Migration: Add last_seen column
-- Description: Track when session was last active

ALTER TABLE sessions 
ADD COLUMN last_seen TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_sessions_last_seen 
ON sessions(last_seen DESC);

COMMENT ON COLUMN sessions.last_seen IS 'Last time session was active';
```

### 3. Escrever o SQL de DOWN (002_add_last_seen.down.sql)

```sql
-- Migration Rollback: Remove last_seen column

DROP INDEX IF EXISTS idx_sessions_last_seen;

ALTER TABLE sessions 
DROP COLUMN IF EXISTS last_seen;
```

### 4. Reiniciar o servidor

```bash
go run cmd/server/main.go
```

A migração será aplicada automaticamente! 🎉

## Verificar Status das Migrações

Conecte ao PostgreSQL e consulte:

```sql
-- Ver todas as migrações aplicadas
SELECT * FROM schema_migrations ORDER BY version;

-- Ver última migração
SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 1;
```

## Reverter uma Migração (Rollback)

⚠️ **CUIDADO**: Rollback pode causar perda de dados!

Para reverter manualmente uma migração específica, você pode usar o código:

```go
package main

import (
    "context"
    "zpwoot/internal/config"
    "zpwoot/internal/db"
    "zpwoot/pkg/logger"
)

func main() {
    logger.InitLogger()
    config.Load()
    db.InitDB()
    
    // Reverter migração versão 2
    err := db.RollbackMigration(context.Background(), 2)
    if err != nil {
        panic(err)
    }
}
```

## Boas Práticas

### ✅ DO

- **Sempre** criar arquivos `.up.sql` e `.down.sql`
- Usar transações quando possível
- Adicionar comentários explicativos
- Testar rollback antes de aplicar em produção
- Fazer backup antes de migrações complexas
- Usar `IF EXISTS` e `IF NOT EXISTS` para idempotência

### ❌ DON'T

- Nunca editar migrações já aplicadas em produção
- Não deletar arquivos de migração
- Evitar migrações que bloqueiam tabelas por muito tempo
- Não fazer migrações destrutivas sem backup

## Troubleshooting

### Erro: "migration already applied"

Isso é normal! O sistema pula migrações já aplicadas automaticamente.

### Erro: "failed to execute migration SQL"

1. Verifique a sintaxe SQL
2. Verifique se a tabela/coluna já existe
3. Verifique permissões do usuário do banco
4. Veja os logs detalhados

### Forçar re-execução de uma migração

```sql
-- CUIDADO: Só faça isso se souber o que está fazendo!
DELETE FROM schema_migrations WHERE version = 1;
```

Depois reinicie o servidor.

## Exemplo Completo

### Adicionar suporte a múltiplos webhooks

**003_add_webhook_retries.up.sql**:
```sql
ALTER TABLE sessions 
ADD COLUMN webhook_retry_count INTEGER DEFAULT 0,
ADD COLUMN webhook_last_error TEXT,
ADD COLUMN webhook_last_retry TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_sessions_webhook_retry 
ON sessions(webhook_retry_count) 
WHERE webhook_retry_count > 0;
```

**003_add_webhook_retries.down.sql**:
```sql
DROP INDEX IF EXISTS idx_sessions_webhook_retry;

ALTER TABLE sessions 
DROP COLUMN IF EXISTS webhook_last_retry,
DROP COLUMN IF EXISTS webhook_last_error,
DROP COLUMN IF EXISTS webhook_retry_count;
```

## Logs

As migrações geram logs detalhados:

```
INFO Starting database migrations...
INFO Applying migration version=1 name=create_sessions
INFO Migration applied successfully version=1 name=create_sessions
INFO All migrations completed successfully
```

## Referências

- [PostgreSQL ALTER TABLE](https://www.postgresql.org/docs/current/sql-altertable.html)
- [PostgreSQL CREATE INDEX](https://www.postgresql.org/docs/current/sql-createindex.html)
- [Go embed package](https://pkg.go.dev/embed)

