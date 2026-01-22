-- ============================================================
-- Migration: 007_content_tables.sql
-- Description: Content tables (news, notifications, files, chat)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- NEWS_CATEGORIES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS news_categories (
    id              SERIAL PRIMARY KEY,
    parent_id       INTEGER REFERENCES news_categories(id) ON DELETE SET NULL,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) UNIQUE NOT NULL,
    description     TEXT,
    icon            VARCHAR(50),
    sort_order      INTEGER DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    deleted_date    TIMESTAMPTZ
);

CREATE INDEX idx_news_categories_code ON news_categories(code);
CREATE INDEX idx_news_categories_parent_id ON news_categories(parent_id);

SELECT create_audit_triggers('news_categories');

-- ============================================================
-- NEWS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS news (
    id                  SERIAL PRIMARY KEY,
    category_id         INTEGER REFERENCES news_categories(id) ON DELETE SET NULL,
    system_id           INTEGER REFERENCES systems(id) ON DELETE SET NULL,
    organization_id     INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    title               VARCHAR(500) NOT NULL,
    slug                VARCHAR(500) UNIQUE,
    summary             TEXT,
    content             TEXT NOT NULL,
    cover_image_url     VARCHAR(500),
    author_id           INTEGER REFERENCES users(id) ON DELETE SET NULL,
    author_name         VARCHAR(200),
    status              VARCHAR(50) DEFAULT 'draft',
    is_featured         BOOLEAN DEFAULT FALSE,
    is_pinned           BOOLEAN DEFAULT FALSE,
    view_count          INTEGER DEFAULT 0,
    like_count          INTEGER DEFAULT 0,
    comment_count       INTEGER DEFAULT 0,
    share_count         INTEGER DEFAULT 0,
    published_at        TIMESTAMPTZ,
    expires_at          TIMESTAMPTZ,
    tags                TEXT[],
    seo_title           VARCHAR(200),
    seo_description     VARCHAR(500),
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ,
    CONSTRAINT chk_news_status CHECK (status IN ('draft', 'pending', 'published', 'archived', 'rejected'))
);

CREATE INDEX idx_news_category_id ON news(category_id);
CREATE INDEX idx_news_system_id ON news(system_id);
CREATE INDEX idx_news_organization_id ON news(organization_id);
CREATE INDEX idx_news_author_id ON news(author_id);
CREATE INDEX idx_news_status ON news(status);
CREATE INDEX idx_news_published_at ON news(published_at);
CREATE INDEX idx_news_is_featured ON news(is_featured);
CREATE INDEX idx_news_slug ON news(slug);
CREATE INDEX idx_news_tags ON news USING GIN(tags);

SELECT create_audit_triggers('news');

-- ============================================================
-- NOTIFICATIONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS notifications (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER REFERENCES users(id) ON DELETE CASCADE,
    system_id           INTEGER REFERENCES systems(id) ON DELETE SET NULL,
    type                VARCHAR(50) NOT NULL,
    title               VARCHAR(255) NOT NULL,
    message             TEXT NOT NULL,
    data                JSONB DEFAULT '{}',
    action_url          VARCHAR(500),
    icon_url            VARCHAR(500),
    priority            VARCHAR(20) DEFAULT 'normal',
    channel             VARCHAR(50) DEFAULT 'push',
    is_read             BOOLEAN DEFAULT FALSE,
    read_at             TIMESTAMPTZ,
    is_sent             BOOLEAN DEFAULT FALSE,
    sent_at             TIMESTAMPTZ,
    scheduled_at        TIMESTAMPTZ,
    expires_at          TIMESTAMPTZ,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_notification_priority CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    CONSTRAINT chk_notification_channel CHECK (channel IN ('push', 'email', 'sms', 'in_app'))
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_system_id ON notifications(system_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_date ON notifications(created_date);

SELECT create_audit_triggers('notifications');

-- ============================================================
-- FILES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS files (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER REFERENCES users(id) ON DELETE SET NULL,
    organization_id     INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    original_name       VARCHAR(500) NOT NULL,
    stored_name         VARCHAR(500) NOT NULL,
    mime_type           VARCHAR(100),
    file_size           BIGINT,
    storage_path        VARCHAR(1000) NOT NULL,
    storage_provider    VARCHAR(50) DEFAULT 'local',
    public_url          VARCHAR(1000),
    thumbnail_url       VARCHAR(1000),
    checksum            VARCHAR(100),
    metadata            JSONB DEFAULT '{}',
    is_public           BOOLEAN DEFAULT FALSE,
    download_count      INTEGER DEFAULT 0,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ
);

CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_organization_id ON files(organization_id);
CREATE INDEX idx_files_mime_type ON files(mime_type);
CREATE INDEX idx_files_storage_provider ON files(storage_provider);

SELECT create_audit_triggers('files');

-- ============================================================
-- CHAT_ROOMS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS chat_rooms (
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR(200),
    type                VARCHAR(50) NOT NULL DEFAULT 'direct',
    organization_id     INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    created_by          INTEGER REFERENCES users(id) ON DELETE SET NULL,
    avatar_url          VARCHAR(500),
    description         TEXT,
    settings            JSONB DEFAULT '{}',
    last_message_at     TIMESTAMPTZ,
    is_active           BOOLEAN DEFAULT TRUE,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    deleted_date        TIMESTAMPTZ,
    CONSTRAINT chk_chat_room_type CHECK (type IN ('direct', 'group', 'channel', 'support'))
);

CREATE INDEX idx_chat_rooms_organization_id ON chat_rooms(organization_id);
CREATE INDEX idx_chat_rooms_type ON chat_rooms(type);
CREATE INDEX idx_chat_rooms_last_message_at ON chat_rooms(last_message_at);

SELECT create_audit_triggers('chat_rooms');

-- ============================================================
-- CHAT_ROOM_MEMBERS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS chat_room_members (
    id                  SERIAL PRIMARY KEY,
    room_id             INTEGER NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id             INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role                VARCHAR(50) DEFAULT 'member',
    nickname            VARCHAR(100),
    last_read_at        TIMESTAMPTZ,
    notifications_muted BOOLEAN DEFAULT FALSE,
    is_active           BOOLEAN DEFAULT TRUE,
    joined_at           TIMESTAMPTZ DEFAULT NOW(),
    left_at             TIMESTAMPTZ,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(room_id, user_id),
    CONSTRAINT chk_member_role CHECK (role IN ('owner', 'admin', 'moderator', 'member'))
);

CREATE INDEX idx_chat_room_members_room_id ON chat_room_members(room_id);
CREATE INDEX idx_chat_room_members_user_id ON chat_room_members(user_id);

SELECT create_audit_triggers('chat_room_members');

-- ============================================================
-- CHAT_MESSAGES TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS chat_messages (
    id                  SERIAL PRIMARY KEY,
    room_id             INTEGER NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    sender_id           INTEGER REFERENCES users(id) ON DELETE SET NULL,
    parent_id           INTEGER REFERENCES chat_messages(id) ON DELETE SET NULL,
    message_type        VARCHAR(50) DEFAULT 'text',
    content             TEXT,
    file_url            VARCHAR(1000),
    file_name           VARCHAR(500),
    file_size           BIGINT,
    metadata            JSONB DEFAULT '{}',
    is_edited           BOOLEAN DEFAULT FALSE,
    edited_at           TIMESTAMPTZ,
    is_deleted          BOOLEAN DEFAULT FALSE,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_message_type CHECK (message_type IN ('text', 'image', 'video', 'audio', 'file', 'location', 'system'))
);

CREATE INDEX idx_chat_messages_room_id ON chat_messages(room_id);
CREATE INDEX idx_chat_messages_sender_id ON chat_messages(sender_id);
CREATE INDEX idx_chat_messages_parent_id ON chat_messages(parent_id);
CREATE INDEX idx_chat_messages_created_date ON chat_messages(created_date);

SELECT create_audit_triggers('chat_messages');
