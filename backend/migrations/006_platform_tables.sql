-- ============================================================
-- Migration: 006_platform_tables.sql
-- Description: Platform tables (app icons, vehicles, etc.)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- APP_ICONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS app_icons (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER NOT NULL REFERENCES systems(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) NOT NULL,
    icon_url        VARCHAR(500),
    deep_link       VARCHAR(500),
    web_url         VARCHAR(500),
    description     TEXT,
    category        VARCHAR(50),
    sort_order      INTEGER DEFAULT 0,
    is_featured     BOOLEAN DEFAULT FALSE,
    is_active       BOOLEAN DEFAULT TRUE,
    start_date      TIMESTAMPTZ,
    end_date        TIMESTAMPTZ,
    target_audience JSONB DEFAULT '{}',
    analytics       JSONB DEFAULT '{}',
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ,
    UNIQUE(system_id, code)
);

CREATE INDEX idx_app_icons_system_id ON app_icons(system_id);
CREATE INDEX idx_app_icons_code ON app_icons(code);
CREATE INDEX idx_app_icons_category ON app_icons(category);
CREATE INDEX idx_app_icons_is_active ON app_icons(is_active);
CREATE INDEX idx_app_icons_is_featured ON app_icons(is_featured);

SELECT create_audit_triggers('app_icons');

-- ============================================================
-- VEHICLES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS vehicles (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER REFERENCES users(id) ON DELETE SET NULL,
    organization_id     INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    plate_number        VARCHAR(20) NOT NULL,
    plate_type          VARCHAR(20),
    vehicle_type        VARCHAR(50),
    brand               VARCHAR(100),
    model               VARCHAR(100),
    color               VARCHAR(50),
    year                INTEGER,
    vin_number          VARCHAR(50),
    engine_number       VARCHAR(50),
    registration_date   DATE,
    inspection_date     DATE,
    insurance_date      DATE,
    photo_url           VARCHAR(500),
    is_verified         BOOLEAN DEFAULT FALSE,
    verified_at         TIMESTAMPTZ,
    verified_by         INTEGER REFERENCES users(id),
    is_active           BOOLEAN DEFAULT TRUE,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ
);

CREATE INDEX idx_vehicles_user_id ON vehicles(user_id);
CREATE INDEX idx_vehicles_organization_id ON vehicles(organization_id);
CREATE INDEX idx_vehicles_plate_number ON vehicles(plate_number);
CREATE INDEX idx_vehicles_is_active ON vehicles(is_active);

SELECT create_audit_triggers('vehicles');

-- ============================================================
-- USER_DEVICES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS user_devices (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id           VARCHAR(255) NOT NULL,
    device_type         VARCHAR(50),
    device_name         VARCHAR(150),
    device_model        VARCHAR(150),
    os_name             VARCHAR(50),
    os_version          VARCHAR(50),
    app_version         VARCHAR(50),
    push_token          TEXT,
    push_provider       VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    last_used_at        TIMESTAMPTZ,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, device_id)
);

CREATE INDEX idx_user_devices_user_id ON user_devices(user_id);
CREATE INDEX idx_user_devices_device_id ON user_devices(device_id);
CREATE INDEX idx_user_devices_is_active ON user_devices(is_active);

SELECT create_audit_triggers('user_devices');

-- ============================================================
-- USER_SETTINGS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS user_settings (
    id                      SERIAL PRIMARY KEY,
    user_id                 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_push       BOOLEAN DEFAULT TRUE,
    notification_email      BOOLEAN DEFAULT TRUE,
    notification_sms        BOOLEAN DEFAULT FALSE,
    theme                   VARCHAR(20) DEFAULT 'system',
    language                VARCHAR(10) DEFAULT 'mn',
    privacy_profile_public  BOOLEAN DEFAULT FALSE,
    privacy_show_email      BOOLEAN DEFAULT FALSE,
    privacy_show_phone      BOOLEAN DEFAULT FALSE,
    custom_settings         JSONB DEFAULT '{}',
    created_date            TIMESTAMPTZ DEFAULT NOW(),
    updated_date            TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id)
);

CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);

SELECT create_audit_triggers('user_settings');
