-- ============================================================
-- Migration: 009_indexes.sql
-- Description: Additional performance indexes and full-text search
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- FULL-TEXT SEARCH INDEXES
-- ============================================================

-- News full-text search
CREATE INDEX idx_news_fts ON news USING GIN(
    to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(summary, '') || ' ' || coalesce(content, ''))
);

-- Users full-text search
CREATE INDEX idx_users_fts ON users USING GIN(
    to_tsvector('simple', coalesce(first_name, '') || ' ' || coalesce(last_name, '') || ' ' || coalesce(email, ''))
);

-- Organizations full-text search
CREATE INDEX idx_organizations_fts ON organizations USING GIN(
    to_tsvector('simple', coalesce(name, '') || ' ' || coalesce(description, ''))
);

-- ============================================================
-- COMPOSITE INDEXES FOR COMMON QUERIES
-- ============================================================

-- User lookup with status
CREATE INDEX idx_users_email_status ON users(email, status) WHERE deleted_date IS NULL;

-- Active sessions lookup
CREATE INDEX idx_user_sessions_active ON user_sessions(user_id, is_active, expires_at)
    WHERE is_active = TRUE AND revoked_at IS NULL;

-- User roles with organization
CREATE INDEX idx_user_roles_active ON user_roles(user_id, role_id, organization_id)
    WHERE is_active = TRUE;

-- News listing
CREATE INDEX idx_news_listing ON news(status, published_at DESC, is_featured, is_pinned)
    WHERE deleted_date IS NULL;

-- Notifications unread
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read, created_date DESC)
    WHERE is_read = FALSE;

-- Login history recent
CREATE INDEX idx_login_history_recent ON login_history(user_id, created_date DESC);

-- ============================================================
-- PARTIAL INDEXES FOR SOFT DELETE PATTERN
-- ============================================================

CREATE INDEX idx_users_active ON users(id) WHERE deleted_date IS NULL;
CREATE INDEX idx_organizations_active ON organizations(id) WHERE deleted_date IS NULL;
CREATE INDEX idx_roles_active ON roles(id) WHERE deleted_date IS NULL;
CREATE INDEX idx_modules_active ON modules(id) WHERE deleted_date IS NULL;
CREATE INDEX idx_permissions_active ON permissions(id) WHERE deleted_date IS NULL;
CREATE INDEX idx_news_active ON news(id) WHERE deleted_date IS NULL;
CREATE INDEX idx_files_active ON files(id) WHERE deleted_date IS NULL;

-- ============================================================
-- JSONB INDEXES
-- ============================================================

-- Organization settings
CREATE INDEX idx_organizations_settings ON organizations USING GIN(settings);

-- User device info
CREATE INDEX idx_user_sessions_device ON user_sessions USING GIN(device_info);

-- News metadata
CREATE INDEX idx_news_seo ON news USING GIN(to_tsvector('simple', coalesce(seo_title, '') || ' ' || coalesce(seo_description, '')));

-- Audit logs values
CREATE INDEX idx_audit_logs_old_values ON audit_logs USING GIN(old_values);
CREATE INDEX idx_audit_logs_new_values ON audit_logs USING GIN(new_values);

-- ============================================================
-- UNIQUE CONSTRAINT INDEXES
-- ============================================================

-- Case-insensitive email uniqueness
CREATE UNIQUE INDEX idx_users_email_lower ON users(LOWER(email)) WHERE deleted_date IS NULL;

-- Case-insensitive organization code
CREATE UNIQUE INDEX idx_organizations_code_lower ON organizations(LOWER(code)) WHERE deleted_date IS NULL;

-- ============================================================
-- STATISTICS
-- ============================================================

-- Analyze tables for query optimization
ANALYZE systems;
ANALYZE actions;
ANALYZE modules;
ANALYZE permissions;
ANALYZE roles;
ANALYZE role_permissions;
ANALYZE menus;
ANALYZE users;
ANALYZE citizens;
ANALYZE user_roles;
ANALYZE user_credentials;
ANALYZE user_sessions;
ANALYZE organizations;
ANALYZE organization_members;
