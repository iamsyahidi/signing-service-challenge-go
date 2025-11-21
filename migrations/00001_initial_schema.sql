-- +goose Up
-- Create the signature_devices table
CREATE TABLE IF NOT EXISTS signature_devices (
    id VARCHAR(255) PRIMARY KEY,
    label VARCHAR(255),
    algorithm VARCHAR(10) NOT NULL,
    signature_counter INTEGER NOT NULL DEFAULT 0,
    last_signature TEXT,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    key_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create an index on the id for faster lookups
CREATE INDEX IF NOT EXISTS idx_signature_devices_id ON signature_devices(id);

-- +goose Down
-- Drop the index
DROP INDEX IF EXISTS idx_signature_devices_id;

-- Drop the signature_devices table
DROP TABLE IF EXISTS signature_devices;
