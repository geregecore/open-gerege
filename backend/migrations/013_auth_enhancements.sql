-- Migration: 013_auth_enhancements.sql
-- Description: Authentication system enhancements
-- - Email verification tokens
-- - Password reset tokens
-- - Refresh tokens
-- - User email verification columns
-- - Backup code salt column
-- Author: Authentication System Refactoring
-- Date: 2026-01-22

-- ============================================================
-- EMAIL VERIFICATION TOKENS
-- ============================================================
-- Хэрэглэгч email баталгаажуулах токенуудыг хадгална.
-- Шинэ бүртгүүлсэн хэрэглэгчдэд 24 цагийн хугацаатай токен илгээгдэнэ.

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(255) UNIQUE NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,

    -- Audit columns (matching existing pattern)
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER,
    created_org_id  INTEGER,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER,
    updated_org_id  INTEGER,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER,
    deleted_org_id  INTEGER
);

-- Index for faster token lookup
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_token ON email_verification_tokens(token);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_expires_at ON email_verification_tokens(expires_at);

COMMENT ON TABLE email_verification_tokens IS 'Email verification tokens for new user registration';
COMMENT ON COLUMN email_verification_tokens.token IS 'Unique verification token (base64 URL encoded)';
COMMENT ON COLUMN email_verification_tokens.expires_at IS 'Token expiration time (typically 24 hours)';
COMMENT ON COLUMN email_verification_tokens.used_at IS 'When token was used (NULL if not yet used)';

-- ============================================================
-- PASSWORD RESET TOKENS
-- ============================================================
-- Хэрэглэгч нууц үг сэргээх токенуудыг хадгална.
-- 1 цагийн хугацаатай токен илгээгдэнэ.

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(255) UNIQUE NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,

    -- Audit columns
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER,
    created_org_id  INTEGER,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER,
    updated_org_id  INTEGER,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER,
    deleted_org_id  INTEGER
);

-- Index for faster token lookup
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);

COMMENT ON TABLE password_reset_tokens IS 'Password reset tokens for forgotten password flow';
COMMENT ON COLUMN password_reset_tokens.token IS 'Unique reset token (base64 URL encoded)';
COMMENT ON COLUMN password_reset_tokens.expires_at IS 'Token expiration time (typically 1 hour)';
COMMENT ON COLUMN password_reset_tokens.used_at IS 'When token was used (NULL if not yet used)';

-- ============================================================
-- REFRESH TOKENS
-- ============================================================
-- Refresh token rotation-д ашиглагдана.
-- Access token: 15 минут, Refresh token: 7 хоног

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) UNIQUE NOT NULL,
    session_id      VARCHAR(255) NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,

    -- Audit columns
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER,
    created_org_id  INTEGER,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER,
    updated_org_id  INTEGER,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER,
    deleted_org_id  INTEGER
);

-- Index for faster token lookup
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_session_id ON refresh_tokens(session_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

COMMENT ON TABLE refresh_tokens IS 'Refresh tokens for token rotation';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'Hashed refresh token (never store plaintext)';
COMMENT ON COLUMN refresh_tokens.session_id IS 'Associated session ID';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'When token was revoked (NULL if active)';

-- ============================================================
-- USERS TABLE UPDATES
-- ============================================================
-- Email баталгаажуулалтын талбарууд нэмэх

-- Add email_verified column if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'email_verified'
    ) THEN
        ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;
    END IF;
END $$;

-- Add email_verified_at column if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'email_verified_at'
    ) THEN
        ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ;
    END IF;
END $$;

COMMENT ON COLUMN users.email_verified IS 'Whether user email has been verified';
COMMENT ON COLUMN users.email_verified_at IS 'When email was verified';

-- ============================================================
-- BACKUP CODE SALT COLUMN
-- ============================================================
-- Backup code-уудад random salt нэмэх (security fix)

-- Add salt column to user_mfa_backup_codes if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'user_mfa_backup_codes' AND column_name = 'salt'
    ) THEN
        ALTER TABLE user_mfa_backup_codes ADD COLUMN salt VARCHAR(64);
    END IF;
END $$;

COMMENT ON COLUMN user_mfa_backup_codes.salt IS 'Random salt for backup code hashing (base64 encoded)';

-- ============================================================
-- CLEANUP OLD TOKENS (scheduled job suggestion)
-- ============================================================
-- Note: Consider setting up a scheduled job to clean up expired tokens:
-- DELETE FROM email_verification_tokens WHERE expires_at < NOW() - INTERVAL '7 days';
-- DELETE FROM password_reset_tokens WHERE expires_at < NOW() - INTERVAL '7 days';
-- DELETE FROM refresh_tokens WHERE expires_at < NOW() - INTERVAL '7 days';
