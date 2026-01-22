-- ============================================================
-- Migration: 001_extensions.sql
-- Description: PostgreSQL extensions and utility functions
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

-- Create schema
CREATE SCHEMA IF NOT EXISTS template_backend;
SET search_path TO template_backend, public;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- TIMESTAMP TRIGGER FUNCTIONS
-- ============================================================

-- Function: Auto-set timestamps on INSERT
CREATE OR REPLACE FUNCTION set_timestamps_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_date := COALESCE(NEW.created_date, NOW());
    NEW.updated_date := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function: Auto-update timestamp on UPDATE
CREATE OR REPLACE FUNCTION set_updated_date_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_date := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- HELPER FUNCTION: Create audit triggers for a table
-- ============================================================

CREATE OR REPLACE FUNCTION create_audit_triggers(table_name TEXT)
RETURNS VOID AS $$
BEGIN
    -- Insert trigger
    EXECUTE format('
        DROP TRIGGER IF EXISTS %I ON %I;
        CREATE TRIGGER %I
        BEFORE INSERT ON %I
        FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
    ',
    table_name || '_insert_timestamp',
    table_name,
    table_name || '_insert_timestamp',
    table_name);

    -- Update trigger
    EXECUTE format('
        DROP TRIGGER IF EXISTS %I ON %I;
        CREATE TRIGGER %I
        BEFORE UPDATE ON %I
        FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();
    ',
    table_name || '_update_timestamp',
    table_name,
    table_name || '_update_timestamp',
    table_name);
END;
$$ LANGUAGE plpgsql;
