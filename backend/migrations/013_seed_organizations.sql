-- ============================================================
-- Migration: 013_seed_organizations.sql
-- Description: Seed data for organization types and organizations
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- ORGANIZATION TYPES
-- ============================================================

INSERT INTO organization_types (name, code, description) VALUES
    ('Төрийн байгууллага', 'GOVERNMENT', 'Төрийн болон төрийн өмчит байгууллагууд'),
    ('Хувийн компани', 'PRIVATE_COMPANY', 'Хувийн хэвшлийн компани'),
    ('ХХК', 'LLC', 'Хязгаарлагдмал хариуцлагатай компани'),
    ('ХК', 'JSC', 'Хувьцаат компани'),
    ('НӨҮГ', 'NGO', 'Төрийн бус байгууллага'),
    ('Банк', 'BANK', 'Банк, санхүүгийн байгууллага'),
    ('Сургууль', 'SCHOOL', 'Боловсролын байгууллага'),
    ('Эмнэлэг', 'HOSPITAL', 'Эрүүл мэндийн байгууллага')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_date = NOW();

-- ============================================================
-- SAMPLE ORGANIZATIONS
-- ============================================================

DO $$
DECLARE
    gov_type_id INTEGER;
    private_type_id INTEGER;
    llc_type_id INTEGER;
    bank_type_id INTEGER;

    -- Organization IDs
    gerege_hq_id INTEGER;
BEGIN
    -- Get type IDs
    SELECT id INTO gov_type_id FROM organization_types WHERE code = 'GOVERNMENT';
    SELECT id INTO private_type_id FROM organization_types WHERE code = 'PRIVATE_COMPANY';
    SELECT id INTO llc_type_id FROM organization_types WHERE code = 'LLC';
    SELECT id INTO bank_type_id FROM organization_types WHERE code = 'BANK';

    -- ========================================
    -- GEREGE HEADQUARTERS
    -- ========================================

    INSERT INTO organizations (type_id, name, code, register_number, description, email, phone, address, aimagcity_code)
    VALUES (
        llc_type_id,
        'Gerege Systems ХХК',
        'GEREGE_HQ',
        '5012345',
        'Gerege системийн үндсэн компани',
        'info@gerege.mn',
        '+976 7011 1234',
        'Улаанбаатар хот, Сүхбаатар дүүрэг, 1-р хороо',
        '01'
    )
    ON CONFLICT (code) DO UPDATE SET
        name = EXCLUDED.name,
        description = EXCLUDED.description,
        updated_date = NOW()
    RETURNING id INTO gerege_hq_id;

    -- ========================================
    -- SAMPLE PARTNER ORGANIZATIONS
    -- ========================================

    -- Government partner
    INSERT INTO organizations (type_id, name, code, description, email, phone, address, aimagcity_code)
    VALUES (
        gov_type_id,
        'Улаанбаатар хотын захиргаа',
        'UB_GOV',
        'Улаанбаатар хотын засаг даргын тамгын газар',
        'info@ulaanbaatar.mn',
        '+976 7026 2626',
        'Улаанбаатар хот, Чингэлтэй дүүрэг',
        '01'
    )
    ON CONFLICT (code) DO UPDATE SET
        name = EXCLUDED.name,
        updated_date = NOW();

    -- Bank partner
    INSERT INTO organizations (type_id, name, code, description, email, phone, address, aimagcity_code)
    VALUES (
        bank_type_id,
        'Хаан банк',
        'KHAN_BANK',
        'Хаан банк',
        'info@khanbank.com',
        '+976 1800 1917',
        'Улаанбаатар хот, Сүхбаатар дүүрэг',
        '01'
    )
    ON CONFLICT (code) DO UPDATE SET
        name = EXCLUDED.name,
        updated_date = NOW();

    -- Private company partner
    INSERT INTO organizations (type_id, name, code, description, email, phone, address, aimagcity_code)
    VALUES (
        private_type_id,
        'Ти Ди Би Монгол ХХК',
        'TDB_MONGOL',
        'Худалдаа, аялал жуулчлал',
        'info@tdbmongol.mn',
        '+976 7011 0000',
        'Улаанбаатар хот',
        '01'
    )
    ON CONFLICT (code) DO UPDATE SET
        name = EXCLUDED.name,
        updated_date = NOW();

END $$;

-- ============================================================
-- NEWS CATEGORIES
-- ============================================================

INSERT INTO news_categories (name, code, description, icon, sort_order) VALUES
    ('Мэдээ', 'NEWS', 'Ерөнхий мэдээ мэдээлэл', 'newspaper', 1),
    ('Зарлал', 'ANNOUNCEMENT', 'Зарлал мэдэгдэл', 'megaphone', 2),
    ('Үйл явдал', 'EVENT', 'Арга хэмжээ, үйл явдал', 'calendar', 3),
    ('Урамшуулал', 'PROMOTION', 'Урамшуулал, хөнгөлөлт', 'gift', 4),
    ('Хууль эрх зүй', 'LEGAL', 'Хууль, журам, дүрэм', 'scale', 5),
    ('Техникийн', 'TECHNICAL', 'Техникийн мэдээлэл, заавар', 'tool', 6)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    updated_date = NOW();
