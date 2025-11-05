-- Migration: Create sessions table
-- Description: Creates the main sessions table for storing WhatsApp session data
-- Author: zpwoot
-- Date: 2025-01-05

-- Enable UUID extension for generating random UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    -- Primary identifier - automatically generated UUID
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,

    -- Session identification
    name TEXT NOT NULL,
    device_jid TEXT,

    -- Session state
    status TEXT NOT NULL DEFAULT 'disconnected',
    connected BOOLEAN DEFAULT FALSE,

    -- WhatsApp data
    qr_code TEXT,

    -- Configuration (JSON)
    proxy_config JSONB DEFAULT NULL,
    webhook_config JSONB DEFAULT NULL,

    -- Authentication - API key for session access
    apikey TEXT DEFAULT NULL,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_sessions_device_jid ON sessions(device_jid);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_apikey ON sessions(apikey);
CREATE INDEX IF NOT EXISTS idx_sessions_connected ON sessions(connected);
CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sessions_proxy_enabled ON sessions ((proxy_config->>'enabled'));
CREATE INDEX IF NOT EXISTS idx_sessions_webhook_enabled ON sessions ((webhook_config->>'enabled'));

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update updated_at on row update
CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to table and columns for documentation
COMMENT ON TABLE sessions IS 'Stores WhatsApp session information and configuration';
COMMENT ON COLUMN sessions.id IS 'Unique session identifier (UUID)';
COMMENT ON COLUMN sessions.name IS 'Human-readable session name';
COMMENT ON COLUMN sessions.device_jid IS 'WhatsApp device JID (after pairing)';
COMMENT ON COLUMN sessions.status IS 'Session status: disconnected, connecting, connected, pairing, failed, logged_out';
COMMENT ON COLUMN sessions.connected IS 'Quick boolean flag for connection state';
COMMENT ON COLUMN sessions.qr_code IS 'Base64 encoded QR code for pairing (temporary)';
COMMENT ON COLUMN sessions.proxy_config IS 'JSON configuration for proxy: {enabled, protocol, host, port, username, password}';
COMMENT ON COLUMN sessions.webhook_config IS 'JSON configuration for webhook: {enabled, url, events, token}';
COMMENT ON COLUMN sessions.apikey IS 'API key for authenticating requests to this session (optional)';
COMMENT ON COLUMN sessions.created_at IS 'Timestamp when session was created';
COMMENT ON COLUMN sessions.updated_at IS 'Timestamp when session was last updated';

