-- ============================================================
-- Migration: 004_auth_tables.sql
-- Description: Authentication tables (credentials, MFA, sessions, tokens)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- USER_CREDENTIALS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS user_credentials (
    id                      SERIAL PRIMARY KEY,
    user_id                 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_type         VARCHAR(50) NOT NULL DEFAULT 'password',
    password_hash           VARCHAR(500),
    password_salt           VARCHAR(100),
    password_changed_at     TIMESTAMPTZ,
    must_change_password    BOOLEAN DEFAULT FALSE,
    oauth_provider          VARCHAR(50),
    oauth_provider_id       VARCHAR(255),
    oauth_access_token      TEXT,
    oauth_refresh_token     TEXT,
    oauth_token_expires_at  TIMESTAMPTZ,
    created_date            TIMESTAMPTZ DEFAULT NOW(),
    updated_date            TIMESTAMPTZ DEFAULT NOW(),
    deleted_date            TIMESTAMPTZ,
    CONSTRAINT chk_credential_type CHECK (credential_type IN ('password', 'oauth', 'dan', 'certificate'))
);

CREATE INDEX idx_user_credentials_user_id ON user_credentials(user_id);
CREATE INDEX idx_user_credentials_type ON user_credentials(credential_type);
CREATE INDEX idx_user_credentials_oauth ON user_credentials(oauth_provider, oauth_provider_id);

SELECT create_audit_triggers('user_credentials');

-- ============================================================
-- USER_MFA_SETTINGS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS user_mfa_settings (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mfa_enabled         BOOLEAN DEFAULT FALSE,
    mfa_type            VARCHAR(50),
    totp_secret         VARCHAR(100),
    totp_verified       BOOLEAN DEFAULT FALSE,
    totp_verified_at    TIMESTAMPTZ,
    recovery_email      VARCHAR(255),
    recovery_phone      VARCHAR(20),
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id),
    CONSTRAINT chk_mfa_type CHECK (mfa_type IS NULL OR mfa_type IN ('totp', 'sms', 'email', 'hardware_key'))
);

CREATE INDEX idx_user_mfa_settings_user_id ON user_mfa_settings(user_id);

SELECT create_audit_triggers('user_mfa_settings');

-- ============================================================
-- USER_MFA_BACKUP_CODES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS user_mfa_backup_codes (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash       VARCHAR(100) NOT NULL,
    salt            VARCHAR(64) NOT NULL,
    used_at         TIMESTAMPTZ,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_mfa_backup_codes_user_id ON user_mfa_backup_codes(user_id);

SELECT create_audit_triggers('user_mfa_backup_codes');

-- ============================================================
-- USER_SESSIONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS user_sessions (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id          VARCHAR(255) UNIQUE NOT NULL,
    access_token_hash   VARCHAR(255),
    refresh_token_hash  VARCHAR(255),
    ip_address          VARCHAR(45),
    user_agent          TEXT,
    device_info         JSONB,
    location            VARCHAR(255),
    is_active           BOOLEAN DEFAULT TRUE,
    last_activity_at    TIMESTAMPTZ DEFAULT NOW(),
    expires_at          TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ,
    revoke_reason       VARCHAR(255),
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_session_id ON user_sessions(session_id);
CREATE INDEX idx_user_sessions_is_active ON user_sessions(is_active);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

SELECT create_audit_triggers('user_sessions');

-- ============================================================
-- LOGIN_HISTORY TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS login_history (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
    email           VARCHAR(255),
    login_type      VARCHAR(50) NOT NULL,
    status          VARCHAR(50) NOT NULL,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    device_info     JSONB,
    location        VARCHAR(255),
    failure_reason  VARCHAR(255),
    session_id      VARCHAR(255),
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_login_status CHECK (status IN ('success', 'failed', 'blocked', 'mfa_required', 'mfa_failed'))
);

CREATE INDEX idx_login_history_user_id ON login_history(user_id);
CREATE INDEX idx_login_history_email ON login_history(email);
CREATE INDEX idx_login_history_status ON login_history(status);
CREATE INDEX idx_login_history_created_date ON login_history(created_date);
CREATE INDEX idx_login_history_ip_address ON login_history(ip_address);

-- ============================================================
-- EMAIL_VERIFICATION_TOKENS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(255) UNIQUE NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);
CREATE INDEX idx_email_verification_tokens_token ON email_verification_tokens(token);

SELECT create_audit_triggers('email_verification_tokens');

-- ============================================================
-- PASSWORD_RESET_TOKENS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(255) UNIQUE NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token);

SELECT create_audit_triggers('password_reset_tokens');

-- ============================================================
-- REFRESH_TOKENS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) UNIQUE NOT NULL,
    session_id      VARCHAR(255) NOT NULL,
    family_id       VARCHAR(255),
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    replaced_by     INTEGER REFERENCES refresh_tokens(id),
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_session_id ON refresh_tokens(session_id);
CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);

SELECT create_audit_triggers('refresh_tokens');
