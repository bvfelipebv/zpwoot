-- Migration Rollback: Drop sessions table
-- Description: Removes the sessions table and related objects
-- Author: zpmeow
-- Date: 2025-01-05

-- Drop trigger
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes (will be dropped automatically with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_sessions_created_at;
DROP INDEX IF EXISTS idx_sessions_connected;
DROP INDEX IF EXISTS idx_sessions_apikey;
DROP INDEX IF EXISTS idx_sessions_status;
DROP INDEX IF EXISTS idx_sessions_device_jid;

-- Drop table
DROP TABLE IF EXISTS sessions;

-- Note: We don't drop the uuid-ossp extension as it might be used by other tables
-- If you want to drop it, uncomment the line below:
-- DROP EXTENSION IF EXISTS "uuid-ossp";

