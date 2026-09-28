-- Enable UUID generation (run once per database)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,                          -- NULL = active (soft delete)

    -- Identity
    email           VARCHAR(255) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    avatar          VARCHAR(500),

    -- Auth
    password_hash   VARCHAR(255),                         -- NULL for OAuth users
    provider        VARCHAR(50)  NOT NULL DEFAULT 'local',
    provider_id     VARCHAR(255),                         -- OAuth provider's user ID

    -- Profile
    role                  VARCHAR(50)  NOT NULL DEFAULT 'student',
    is_verified           BOOLEAN NOT NULL DEFAULT FALSE,
    is_active             BOOLEAN NOT NULL DEFAULT TRUE,
    preferred_roles       JSONB DEFAULT '[]',
    preferred_locations   JSONB DEFAULT '[]',
    min_salary            INTEGER NOT NULL DEFAULT 0,

    -- Constraints
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_email_not_empty CHECK (email <> '')
);

-- Index on email for fast login lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index on deleted_at to filter soft-deleted users efficiently
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
