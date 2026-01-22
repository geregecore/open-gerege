-- ============================================================
-- Migration: 012_seed_roles.sql
-- Description: Seed data for roles and role_permissions
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- ROLES
-- ============================================================

DO $$
DECLARE
    admin_system_id INTEGER;
    gerege_app_id INTEGER;
    gerege_business_id INTEGER;
    tpay_system_id INTEGER;

    -- Role IDs
    super_admin_id INTEGER;
    admin_id INTEGER;
    moderator_id INTEGER;
    support_id INTEGER;
    app_user_id INTEGER;
    app_premium_id INTEGER;
    biz_owner_id INTEGER;
    biz_manager_id INTEGER;
    biz_employee_id INTEGER;
BEGIN
    -- Get system IDs
    SELECT id INTO admin_system_id FROM systems WHERE code = 'ADMIN';
    SELECT id INTO gerege_app_id FROM systems WHERE code = 'GEREGE_APP';
    SELECT id INTO gerege_business_id FROM systems WHERE code = 'GEREGE_BUSINESS';
    SELECT id INTO tpay_system_id FROM systems WHERE code = 'TPAY';

    -- ========================================
    -- ADMIN SYSTEM ROLES
    -- ========================================

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (admin_system_id, 'Супер Админ', 'SUPER_ADMIN', 'Бүх эрхтэй систем админ', TRUE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO super_admin_id;

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (admin_system_id, 'Админ', 'ADMIN', 'Системийн админ', TRUE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO admin_id;

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (admin_system_id, 'Модератор', 'MODERATOR', 'Контент модератор', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO moderator_id;

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (admin_system_id, 'Дэмжлэг', 'SUPPORT', 'Хэрэглэгчийн дэмжлэг', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO support_id;

    -- ========================================
    -- GEREGE APP ROLES
    -- ========================================

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (gerege_app_id, 'Хэрэглэгч', 'APP_USER', 'Энгийн хэрэглэгч', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO app_user_id;

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (gerege_app_id, 'Премиум хэрэглэгч', 'APP_PREMIUM', 'Премиум эрхтэй хэрэглэгч', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO app_premium_id;

    -- ========================================
    -- GEREGE BUSINESS ROLES
    -- ========================================

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (gerege_business_id, 'Эзэмшигч', 'BIZ_OWNER', 'Байгууллагын эзэмшигч', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO biz_owner_id;

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (gerege_business_id, 'Менежер', 'BIZ_MANAGER', 'Байгууллагын менежер', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO biz_manager_id;

    INSERT INTO roles (system_id, name, code, description, is_system_role) VALUES
        (gerege_business_id, 'Ажилтан', 'BIZ_EMPLOYEE', 'Байгууллагын ажилтан', FALSE)
    ON CONFLICT (system_id, code) DO UPDATE SET name = EXCLUDED.name, updated_date = NOW()
    RETURNING id INTO biz_employee_id;

    -- ========================================
    -- ROLE PERMISSIONS - SUPER_ADMIN (All permissions)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT super_admin_id, p.id
    FROM permissions p
    WHERE p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - ADMIN (Most permissions except system config)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT admin_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    WHERE m.code NOT IN ('SYSTEMS', 'MODULES', 'SCHEDULED_JOBS')
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - MODERATOR (Content only)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT moderator_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    WHERE m.code IN ('NEWS', 'NOTIFICATIONS', 'FILES', 'APP_ICONS', 'CONTENT')
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - SUPPORT (Users read-only + history)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT support_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    JOIN actions a ON p.action_id = a.id
    WHERE (
        (m.code IN ('USERS_LIST', 'CITIZENS', 'SESSIONS', 'LOGIN_HISTORY') AND a.code IN ('VIEW', 'LIST'))
        OR m.code = 'DASHBOARD'
    )
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - APP_USER
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT app_user_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    WHERE m.code IN ('APP_HOME', 'APP_PROFILE', 'APP_VEHICLES', 'APP_NOTIFICATIONS', 'APP_SETTINGS')
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - APP_PREMIUM (All app permissions)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT app_premium_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    JOIN systems s ON m.system_id = s.id
    WHERE s.code = 'GEREGE_APP'
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - BIZ_OWNER (All business permissions)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT biz_owner_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    JOIN systems s ON m.system_id = s.id
    WHERE s.code = 'GEREGE_BUSINESS'
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - BIZ_MANAGER
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT biz_manager_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    JOIN actions a ON p.action_id = a.id
    WHERE m.code IN ('BIZ_DASHBOARD', 'BIZ_EMPLOYEES', 'BIZ_VEHICLES', 'BIZ_REPORTS')
    AND a.code != 'DELETE'
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

    -- ========================================
    -- ROLE PERMISSIONS - BIZ_EMPLOYEE (Read only)
    -- ========================================

    INSERT INTO role_permissions (role_id, permission_id)
    SELECT biz_employee_id, p.id
    FROM permissions p
    JOIN modules m ON p.module_id = m.id
    JOIN actions a ON p.action_id = a.id
    WHERE m.code IN ('BIZ_DASHBOARD', 'BIZ_VEHICLES')
    AND a.code IN ('VIEW', 'LIST')
    AND p.deleted_date IS NULL
    ON CONFLICT (role_id, permission_id) DO NOTHING;

END $$;
