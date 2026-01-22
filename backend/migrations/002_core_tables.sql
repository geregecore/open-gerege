-- ============================================================
-- Migration: 002_core_tables.sql
-- Description: Core system tables (systems, actions, modules, permissions, roles, menus)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- SYSTEMS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS systems (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) UNIQUE NOT NULL,
    description     TEXT,
    base_url        VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ
);

CREATE INDEX idx_systems_code ON systems(code);
CREATE INDEX idx_systems_is_active ON systems(is_active);

SELECT create_audit_triggers('systems');

-- ============================================================
-- ACTIONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS actions (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) UNIQUE NOT NULL,
    description     TEXT,
    http_method     VARCHAR(10),
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ
);

CREATE INDEX idx_actions_code ON actions(code);

SELECT create_audit_triggers('actions');

-- ============================================================
-- MODULES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS modules (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER NOT NULL REFERENCES systems(id) ON DELETE CASCADE,
    parent_id       INTEGER REFERENCES modules(id) ON DELETE SET NULL,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(100) NOT NULL,
    description     TEXT,
    icon            VARCHAR(50),
    sort_order      INTEGER DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ,
    UNIQUE(system_id, code)
);

CREATE INDEX idx_modules_system_id ON modules(system_id);
CREATE INDEX idx_modules_parent_id ON modules(parent_id);
CREATE INDEX idx_modules_code ON modules(code);
CREATE INDEX idx_modules_is_active ON modules(is_active);

SELECT create_audit_triggers('modules');

-- ============================================================
-- PERMISSIONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS permissions (
    id              SERIAL PRIMARY KEY,
    module_id       INTEGER NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    action_id       INTEGER NOT NULL REFERENCES actions(id) ON DELETE CASCADE,
    name            VARCHAR(150) NOT NULL,
    code            VARCHAR(150) NOT NULL,
    description     TEXT,
    resource_path   VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ,
    UNIQUE(module_id, action_id)
);

CREATE INDEX idx_permissions_module_id ON permissions(module_id);
CREATE INDEX idx_permissions_action_id ON permissions(action_id);
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_is_active ON permissions(is_active);

SELECT create_audit_triggers('permissions');

-- ============================================================
-- ROLES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS roles (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER NOT NULL REFERENCES systems(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) NOT NULL,
    description     TEXT,
    is_system_role  BOOLEAN DEFAULT FALSE,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ,
    UNIQUE(system_id, code)
);

CREATE INDEX idx_roles_system_id ON roles(system_id);
CREATE INDEX idx_roles_code ON roles(code);
CREATE INDEX idx_roles_is_active ON roles(is_active);

SELECT create_audit_triggers('roles');

-- ============================================================
-- ROLE_PERMISSIONS TABLE (Many-to-Many)
-- ============================================================

CREATE TABLE IF NOT EXISTS role_permissions (
    id              SERIAL PRIMARY KEY,
    role_id         INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id   INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

SELECT create_audit_triggers('role_permissions');

-- ============================================================
-- MENUS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS menus (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER NOT NULL REFERENCES systems(id) ON DELETE CASCADE,
    parent_id       INTEGER REFERENCES menus(id) ON DELETE SET NULL,
    module_id       INTEGER REFERENCES modules(id) ON DELETE SET NULL,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(100) NOT NULL,
    icon            VARCHAR(50),
    url             VARCHAR(255),
    sort_order      INTEGER DEFAULT 0,
    is_visible      BOOLEAN DEFAULT TRUE,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ,
    UNIQUE(system_id, code)
);

CREATE INDEX idx_menus_system_id ON menus(system_id);
CREATE INDEX idx_menus_parent_id ON menus(parent_id);
CREATE INDEX idx_menus_module_id ON menus(module_id);
CREATE INDEX idx_menus_is_active ON menus(is_active);

SELECT create_audit_triggers('menus');
