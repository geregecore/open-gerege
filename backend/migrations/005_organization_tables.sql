-- ============================================================
-- Migration: 005_organization_tables.sql
-- Description: Organization tables
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- ORGANIZATION_TYPES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS organization_types (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) UNIQUE NOT NULL,
    description     TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ
);

CREATE INDEX idx_organization_types_code ON organization_types(code);

SELECT create_audit_triggers('organization_types');

-- ============================================================
-- ORGANIZATIONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS organizations (
    id                  SERIAL PRIMARY KEY,
    type_id             INTEGER NOT NULL REFERENCES organization_types(id) ON DELETE RESTRICT,
    parent_id           INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    name                VARCHAR(255) NOT NULL,
    code                VARCHAR(100) UNIQUE,
    register_number     VARCHAR(50),
    description         TEXT,
    logo_url            VARCHAR(500),
    website             VARCHAR(255),
    email               VARCHAR(255),
    phone               VARCHAR(50),
    address             TEXT,
    aimagcity_code      VARCHAR(10),
    soum_code           VARCHAR(10),
    latitude            DECIMAL(10, 8),
    longitude           DECIMAL(11, 8),
    settings            JSONB DEFAULT '{}',
    is_active           BOOLEAN DEFAULT TRUE,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ
);

CREATE INDEX idx_organizations_type_id ON organizations(type_id);
CREATE INDEX idx_organizations_parent_id ON organizations(parent_id);
CREATE INDEX idx_organizations_code ON organizations(code);
CREATE INDEX idx_organizations_is_active ON organizations(is_active);
CREATE INDEX idx_organizations_aimagcity ON organizations(aimagcity_code);

SELECT create_audit_triggers('organizations');

-- ============================================================
-- ORGANIZATION_MEMBERS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS organization_members (
    id              SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position        VARCHAR(150),
    department      VARCHAR(150),
    employee_id     VARCHAR(50),
    start_date      DATE,
    end_date        DATE,
    is_primary      BOOLEAN DEFAULT FALSE,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ,
    UNIQUE(organization_id, user_id)
);

CREATE INDEX idx_organization_members_org_id ON organization_members(organization_id);
CREATE INDEX idx_organization_members_user_id ON organization_members(user_id);
CREATE INDEX idx_organization_members_is_active ON organization_members(is_active);

SELECT create_audit_triggers('organization_members');

-- ============================================================
-- Add organization_id FK to user_roles (deferred from 003)
-- ============================================================

ALTER TABLE user_roles
ADD CONSTRAINT fk_user_roles_organization
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;
