-- Inicialização do banco de dados zpmeow
-- Este script é executado automaticamente quando o container PostgreSQL é criado

-- Criar extensões úteis
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Configurar timezone
SET timezone = 'UTC';

-- Criar schema se necessário (opcional)
-- CREATE SCHEMA IF NOT EXISTS zpmeow;

-- O GORM criará as tabelas automaticamente via AutoMigrate
-- Mas você pode adicionar índices adicionais ou configurações aqui se necessário

-- Exemplo de índice adicional (descomente se necessário):
-- CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
-- CREATE INDEX IF NOT EXISTS idx_sessions_jid ON sessions(jid);

