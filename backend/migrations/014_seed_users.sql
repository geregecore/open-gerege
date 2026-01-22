-- ============================================================
-- Migration: 014_seed_users.sql
-- Description: Seed data for users (admin accounts)
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- ADMIN USER ACCOUNTS
-- ============================================================

DO $$
DECLARE
    super_admin_user_id INTEGER;
    admin_user_id INTEGER;

    super_admin_role_id INTEGER;
    admin_role_id INTEGER;

    gerege_org_id INTEGER;
BEGIN
    -- Get role IDs
    SELECT id INTO super_admin_role_id FROM roles WHERE code = 'SUPER_ADMIN';
    SELECT id INTO admin_role_id FROM roles WHERE code = 'ADMIN';

    -- Get Gerege organization ID
    SELECT id INTO gerege_org_id FROM organizations WHERE code = 'GEREGE_HQ';

    -- ========================================
    -- SUPER ADMIN USER
    -- ========================================

    INSERT INTO users (email, phone, first_name, last_name, status, email_verified, email_verified_at)
    VALUES (
        'superadmin@gerege.mn',
        '+976 9911 0001',
        'Super',
        'Admin',
        'active',
        TRUE,
        NOW()
    )
    ON CONFLICT (email) DO UPDATE SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        status = EXCLUDED.status,
        updated_date = NOW()
    RETURNING id INTO super_admin_user_id;

    -- Create credentials for super admin
    -- Password: Gerege@2024! (hashed with Argon2id - this is a placeholder, actual hash should be generated)
    INSERT INTO user_credentials (user_id, credential_type, password_hash, password_salt, password_changed_at)
    VALUES (
        super_admin_user_id,
        'password',
        '$argon2id$v=19$m=65536,t=1,p=4$c29tZXNhbHQ$placeholder_hash_replace_in_production',
        'somesalt_replace_in_production',
        NOW()
    )
    ON CONFLICT DO NOTHING;

    -- Assign super admin role
    INSERT INTO user_roles (user_id, role_id, organization_id, assigned_at)
    VALUES (super_admin_user_id, super_admin_role_id, gerege_org_id, NOW())
    ON CONFLICT (user_id, role_id, organization_id) DO NOTHING;

    -- Create user settings
    INSERT INTO user_settings (user_id, notification_push, notification_email, theme, language)
    VALUES (super_admin_user_id, TRUE, TRUE, 'system', 'mn')
    ON CONFLICT (user_id) DO NOTHING;

    -- ========================================
    -- ADMIN USER
    -- ========================================

    INSERT INTO users (email, phone, first_name, last_name, status, email_verified, email_verified_at)
    VALUES (
        'admin@gerege.mn',
        '+976 9911 0002',
        'System',
        'Admin',
        'active',
        TRUE,
        NOW()
    )
    ON CONFLICT (email) DO UPDATE SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        status = EXCLUDED.status,
        updated_date = NOW()
    RETURNING id INTO admin_user_id;

    -- Create credentials for admin
    INSERT INTO user_credentials (user_id, credential_type, password_hash, password_salt, password_changed_at)
    VALUES (
        admin_user_id,
        'password',
        '$argon2id$v=19$m=65536,t=1,p=4$c29tZXNhbHQ$placeholder_hash_replace_in_production',
        'somesalt_replace_in_production',
        NOW()
    )
    ON CONFLICT DO NOTHING;

    -- Assign admin role
    INSERT INTO user_roles (user_id, role_id, organization_id, assigned_at)
    VALUES (admin_user_id, admin_role_id, gerege_org_id, NOW())
    ON CONFLICT (user_id, role_id, organization_id) DO NOTHING;

    -- Create user settings
    INSERT INTO user_settings (user_id, notification_push, notification_email, theme, language)
    VALUES (admin_user_id, TRUE, TRUE, 'system', 'mn')
    ON CONFLICT (user_id) DO NOTHING;

    -- ========================================
    -- ADD ORGANIZATION MEMBERSHIP
    -- ========================================

    INSERT INTO organization_members (organization_id, user_id, position, department, is_primary, is_active)
    VALUES
        (gerege_org_id, super_admin_user_id, 'Системийн админ', 'IT', TRUE, TRUE),
        (gerege_org_id, admin_user_id, 'Админ', 'IT', TRUE, TRUE)
    ON CONFLICT (organization_id, user_id) DO UPDATE SET
        position = EXCLUDED.position,
        updated_date = NOW();

END $$;

-- ============================================================
-- IMPORTANT NOTE FOR PRODUCTION
-- ============================================================

-- The password hashes above are placeholders!
-- In production, you must:
-- 1. Generate proper Argon2id hashes for passwords
-- 2. Use secure random salts
-- 3. Never commit actual credentials to version control
-- 4. Change default passwords immediately after deployment

-- Example of generating proper hash in Go:
-- hash := argon2.IDKey([]byte("password"), salt, 1, 64*1024, 4, 32)

-- For initial setup, use the application's password reset flow
-- or update the hash directly in the database after generating
-- it with the application's password hashing function.
