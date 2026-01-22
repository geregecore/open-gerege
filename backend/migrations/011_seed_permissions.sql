-- ============================================================
-- Migration: 011_seed_permissions.sql
-- Description: Seed data for permissions
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- PERMISSIONS - Generated for each module + action combination
-- ============================================================

DO $$
DECLARE
    -- Action IDs
    action_view_id INTEGER;
    action_list_id INTEGER;
    action_create_id INTEGER;
    action_update_id INTEGER;
    action_delete_id INTEGER;
    action_export_id INTEGER;
    action_approve_id INTEGER;

    -- Module record
    mod RECORD;
BEGIN
    -- Get action IDs
    SELECT id INTO action_view_id FROM actions WHERE code = 'VIEW';
    SELECT id INTO action_list_id FROM actions WHERE code = 'LIST';
    SELECT id INTO action_create_id FROM actions WHERE code = 'CREATE';
    SELECT id INTO action_update_id FROM actions WHERE code = 'UPDATE';
    SELECT id INTO action_delete_id FROM actions WHERE code = 'DELETE';
    SELECT id INTO action_export_id FROM actions WHERE code = 'EXPORT';
    SELECT id INTO action_approve_id FROM actions WHERE code = 'APPROVE';

    -- Generate permissions for each module
    FOR mod IN SELECT id, code, name FROM modules WHERE deleted_date IS NULL LOOP
        -- VIEW permission
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_view_id, mod.name || ' - Үзэх', mod.code || '_VIEW', '/' || LOWER(mod.code))
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();

        -- LIST permission
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_list_id, mod.name || ' - Жагсаалт', mod.code || '_LIST', '/' || LOWER(mod.code))
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();

        -- CREATE permission
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_create_id, mod.name || ' - Нэмэх', mod.code || '_CREATE', '/' || LOWER(mod.code))
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();

        -- UPDATE permission
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_update_id, mod.name || ' - Засах', mod.code || '_UPDATE', '/' || LOWER(mod.code))
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();

        -- DELETE permission
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_delete_id, mod.name || ' - Устгах', mod.code || '_DELETE', '/' || LOWER(mod.code))
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();
    END LOOP;

    -- Add EXPORT permissions for report modules
    FOR mod IN SELECT id, code, name FROM modules WHERE code LIKE '%REPORT%' OR code LIKE '%LOG%' LOOP
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_export_id, mod.name || ' - Экспорт', mod.code || '_EXPORT', '/' || LOWER(mod.code) || '/export')
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();
    END LOOP;

    -- Add APPROVE permissions for content modules
    FOR mod IN SELECT id, code, name FROM modules WHERE code IN ('NEWS', 'ORGANIZATIONS') LOOP
        INSERT INTO permissions (module_id, action_id, name, code, resource_path)
        VALUES (mod.id, action_approve_id, mod.name || ' - Батлах', mod.code || '_APPROVE', '/' || LOWER(mod.code) || '/approve')
        ON CONFLICT (module_id, action_id) DO UPDATE SET
            name = EXCLUDED.name,
            code = EXCLUDED.code,
            updated_date = NOW();
    END LOOP;

END $$;
