-- ============================================================
-- Migration: 010_seed_core.sql
-- Description: Seed data for core tables (systems, actions, modules)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- SYSTEMS
-- ============================================================

INSERT INTO systems (name, code, description, base_url, is_active) VALUES
    ('Админ систем', 'ADMIN', 'Системийн удирдлагын админ панел', '/admin', TRUE),
    ('Gerege App', 'GEREGE_APP', 'Gerege гар утасны аппликейшн', '/app', TRUE),
    ('Gerege Business', 'GEREGE_BUSINESS', 'Байгууллагуудад зориулсан систем', '/business', TRUE),
    ('TPay', 'TPAY', 'Төлбөрийн систем', '/tpay', TRUE)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    base_url = EXCLUDED.base_url,
    updated_date = NOW();

-- ============================================================
-- ACTIONS (CRUD + Custom)
-- ============================================================

INSERT INTO actions (name, code, description, http_method) VALUES
    ('Үзэх', 'VIEW', 'Мэдээлэл харах', 'GET'),
    ('Жагсаалт', 'LIST', 'Жагсаалт харах', 'GET'),
    ('Нэмэх', 'CREATE', 'Шинээр үүсгэх', 'POST'),
    ('Засах', 'UPDATE', 'Мэдээлэл засах', 'PUT'),
    ('Устгах', 'DELETE', 'Мэдээлэл устгах', 'DELETE'),
    ('Экспорт', 'EXPORT', 'Тайлан татах', 'GET'),
    ('Импорт', 'IMPORT', 'Өгөгдөл импортлох', 'POST'),
    ('Батлах', 'APPROVE', 'Баталгаажуулах', 'POST'),
    ('Татгалзах', 'REJECT', 'Татгалзах', 'POST'),
    ('Хэвлэх', 'PRINT', 'Хэвлэх', 'GET')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    http_method = EXCLUDED.http_method,
    updated_date = NOW();

-- ============================================================
-- MODULES - ADMIN SYSTEM
-- ============================================================

-- Get admin system ID
DO $$
DECLARE
    admin_system_id INTEGER;
    gerege_app_id INTEGER;
    gerege_business_id INTEGER;
    tpay_system_id INTEGER;

    -- Module IDs
    dashboard_id INTEGER;
    users_id INTEGER;
    roles_id INTEGER;
    organizations_id INTEGER;
    content_id INTEGER;
    settings_id INTEGER;
    reports_id INTEGER;
BEGIN
    SELECT id INTO admin_system_id FROM systems WHERE code = 'ADMIN';
    SELECT id INTO gerege_app_id FROM systems WHERE code = 'GEREGE_APP';
    SELECT id INTO gerege_business_id FROM systems WHERE code = 'GEREGE_BUSINESS';
    SELECT id INTO tpay_system_id FROM systems WHERE code = 'TPAY';

    -- ========================================
    -- ADMIN SYSTEM MODULES
    -- ========================================

    -- Dashboard
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Хянах самбар', 'DASHBOARD', 'Системийн ерөнхий мэдээлэл', 'dashboard', 1)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO dashboard_id;

    -- Users Management
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Хэрэглэгчид', 'USERS', 'Хэрэглэгчийн удирдлага', 'users', 2)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO users_id;

    -- Sub-modules for Users
    INSERT INTO modules (system_id, parent_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, users_id, 'Хэрэглэгчийн жагсаалт', 'USERS_LIST', 'Бүх хэрэглэгчид', 'list', 1),
        (admin_system_id, users_id, 'Иргэний мэдээлэл', 'CITIZENS', 'ДАН-аас баталгаажсан иргэд', 'id-card', 2),
        (admin_system_id, users_id, 'Сессионууд', 'SESSIONS', 'Идэвхтэй сессионууд', 'monitor', 3)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- Roles & Permissions
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Эрхүүд', 'ROLES', 'Эрхийн удирдлага', 'shield', 3)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO roles_id;

    INSERT INTO modules (system_id, parent_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, roles_id, 'Роль', 'ROLES_LIST', 'Ролиуд', 'key', 1),
        (admin_system_id, roles_id, 'Зөвшөөрөл', 'PERMISSIONS', 'Зөвшөөрлүүд', 'lock', 2)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- Organizations
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Байгууллага', 'ORGANIZATIONS', 'Байгууллагын удирдлага', 'building', 4)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO organizations_id;

    INSERT INTO modules (system_id, parent_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, organizations_id, 'Байгууллагын жагсаалт', 'ORGANIZATIONS_LIST', 'Бүх байгууллага', 'list', 1),
        (admin_system_id, organizations_id, 'Байгууллагын төрөл', 'ORGANIZATION_TYPES', 'Төрлүүд', 'category', 2)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- Content Management
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Контент', 'CONTENT', 'Контент удирдлага', 'file-text', 5)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO content_id;

    INSERT INTO modules (system_id, parent_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, content_id, 'Мэдээ', 'NEWS', 'Мэдээ мэдээлэл', 'newspaper', 1),
        (admin_system_id, content_id, 'Мэдэгдэл', 'NOTIFICATIONS', 'Мэдэгдлүүд', 'bell', 2),
        (admin_system_id, content_id, 'Файл', 'FILES', 'Файлын сан', 'folder', 3),
        (admin_system_id, content_id, 'Апп icon', 'APP_ICONS', 'Аппликейшн icon', 'grid', 4)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- Settings
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Тохиргоо', 'SETTINGS', 'Системийн тохиргоо', 'settings', 6)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO settings_id;

    INSERT INTO modules (system_id, parent_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, settings_id, 'Системүүд', 'SYSTEMS', 'Системийн тохиргоо', 'server', 1),
        (admin_system_id, settings_id, 'Модулиуд', 'MODULES', 'Модулийн тохиргоо', 'layout', 2),
        (admin_system_id, settings_id, 'Меню', 'MENUS', 'Менюний тохиргоо', 'menu', 3),
        (admin_system_id, settings_id, 'Цагийн ажил', 'SCHEDULED_JOBS', 'Товлосон ажлууд', 'clock', 4)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- Reports
    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, 'Тайлан', 'REPORTS', 'Тайлан статистик', 'bar-chart', 7)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO reports_id;

    INSERT INTO modules (system_id, parent_id, name, code, description, icon, sort_order) VALUES
        (admin_system_id, reports_id, 'Аудит лог', 'AUDIT_LOGS', 'Системийн лог', 'file-search', 1),
        (admin_system_id, reports_id, 'API лог', 'API_LOGS', 'API хүсэлтийн лог', 'terminal', 2),
        (admin_system_id, reports_id, 'Алдааны лог', 'ERROR_LOGS', 'Алдааны бүртгэл', 'alert-triangle', 3),
        (admin_system_id, reports_id, 'Нэвтрэлтийн түүх', 'LOGIN_HISTORY', 'Нэвтрэлтийн түүх', 'log-in', 4)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- ========================================
    -- GEREGE APP MODULES
    -- ========================================

    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (gerege_app_id, 'Нүүр', 'APP_HOME', 'Аппын нүүр хуудас', 'home', 1),
        (gerege_app_id, 'Профайл', 'APP_PROFILE', 'Хэрэглэгчийн профайл', 'user', 2),
        (gerege_app_id, 'Тээврийн хэрэгсэл', 'APP_VEHICLES', 'Миний машин', 'car', 3),
        (gerege_app_id, 'Мэдэгдэл', 'APP_NOTIFICATIONS', 'Мэдэгдлүүд', 'bell', 4),
        (gerege_app_id, 'Тохиргоо', 'APP_SETTINGS', 'Тохиргоо', 'settings', 5)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- ========================================
    -- GEREGE BUSINESS MODULES
    -- ========================================

    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (gerege_business_id, 'Хянах самбар', 'BIZ_DASHBOARD', 'Бизнес хянах самбар', 'dashboard', 1),
        (gerege_business_id, 'Ажилтнууд', 'BIZ_EMPLOYEES', 'Ажилтны удирдлага', 'users', 2),
        (gerege_business_id, 'Тээврийн хэрэгсэл', 'BIZ_VEHICLES', 'Байгууллагын машин', 'truck', 3),
        (gerege_business_id, 'Тайлан', 'BIZ_REPORTS', 'Тайлан', 'file-text', 4),
        (gerege_business_id, 'Тохиргоо', 'BIZ_SETTINGS', 'Байгууллагын тохиргоо', 'settings', 5)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

    -- ========================================
    -- TPAY MODULES
    -- ========================================

    INSERT INTO modules (system_id, name, code, description, icon, sort_order) VALUES
        (tpay_system_id, 'Төлбөр', 'TPAY_PAYMENTS', 'Төлбөр хийх', 'credit-card', 1),
        (tpay_system_id, 'Гүйлгээ', 'TPAY_TRANSACTIONS', 'Гүйлгээний түүх', 'list', 2),
        (tpay_system_id, 'Данс', 'TPAY_ACCOUNTS', 'Дансны мэдээлэл', 'wallet', 3)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW();

END $$;
