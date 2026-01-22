-- ============================================================
-- Migration: 003_user_tables.sql
-- Description: User tables (users, citizens, user_roles)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- CITIZENS TABLE (Иргэний мэдээлэл - ДАН-аас)
-- ============================================================

CREATE TABLE IF NOT EXISTS citizens (
    id                  SERIAL PRIMARY KEY,
    register_number     VARCHAR(20) UNIQUE,
    first_name          VARCHAR(150),
    last_name           VARCHAR(150),
    surname             VARCHAR(150),
    date_of_birth       DATE,
    gender              VARCHAR(10),
    nationality         VARCHAR(50),
    aimagcity_name      VARCHAR(100),
    soum_name           VARCHAR(100),
    bag_name            VARCHAR(100),
    address_detail      TEXT,
    passport_number     VARCHAR(50),
    passport_issue_date DATE,
    passport_expiry_date DATE,
    photo_base64        TEXT,
    dan_verified        BOOLEAN DEFAULT FALSE,
    dan_verified_at     TIMESTAMPTZ,
    dan_raw_data        JSONB,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ
);

CREATE INDEX idx_citizens_register_number ON citizens(register_number);
CREATE INDEX idx_citizens_first_name ON citizens(first_name);
CREATE INDEX idx_citizens_last_name ON citizens(last_name);

SELECT create_audit_triggers('citizens');

-- ============================================================
-- USERS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id                  SERIAL PRIMARY KEY,
    citizen_id          INTEGER REFERENCES citizens(id) ON DELETE SET NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    phone               VARCHAR(20),
    first_name          VARCHAR(150) NOT NULL,
    last_name           VARCHAR(150) NOT NULL,
    avatar_url          VARCHAR(500),
    language            VARCHAR(10) DEFAULT 'mn',
    timezone            VARCHAR(50) DEFAULT 'Asia/Ulaanbaatar',
    status              VARCHAR(50) DEFAULT 'pending_verification',
    email_verified      BOOLEAN DEFAULT FALSE,
    email_verified_at   TIMESTAMPTZ,
    phone_verified      BOOLEAN DEFAULT FALSE,
    phone_verified_at   TIMESTAMPTZ,
    last_login_at       TIMESTAMPTZ,
    last_login_ip       VARCHAR(45),
    failed_login_count  INTEGER DEFAULT 0,
    locked_until        TIMESTAMPTZ,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ,
    CONSTRAINT chk_user_status CHECK (status IN ('pending_verification', 'active', 'inactive', 'suspended', 'locked'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_citizen_id ON users(citizen_id);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_date ON users(deleted_date);

SELECT create_audit_triggers('users');

-- ============================================================
-- USER_ROLES TABLE (Many-to-Many)
-- ============================================================

CREATE TABLE IF NOT EXISTS user_roles (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id         INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    organization_id INTEGER,
    assigned_by     INTEGER REFERENCES users(id) ON DELETE SET NULL,
    assigned_at     TIMESTAMPTZ DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, role_id, organization_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_organization_id ON user_roles(organization_id);
CREATE INDEX idx_user_roles_is_active ON user_roles(is_active);

SELECT create_audit_triggers('user_roles');
