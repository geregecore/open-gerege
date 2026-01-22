-- ============================================================
-- Migration: 008_logging_tables.sql
-- Description: Logging and audit tables
-- Database: gerege_db
-- Schema: template_backend
-- ============================================================

SET search_path TO template_backend, public;

-- ============================================================
-- AUDIT_LOGS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER REFERENCES users(id) ON DELETE SET NULL,
    organization_id     INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    action              VARCHAR(100) NOT NULL,
    entity_type         VARCHAR(100),
    entity_id           INTEGER,
    old_values          JSONB,
    new_values          JSONB,
    ip_address          VARCHAR(45),
    user_agent          TEXT,
    session_id          VARCHAR(255),
    request_id          VARCHAR(255),
    status              VARCHAR(50) DEFAULT 'success',
    error_message       TEXT,
    duration_ms         INTEGER,
    created_date        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_organization_id ON audit_logs(organization_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_date ON audit_logs(created_date);
CREATE INDEX idx_audit_logs_request_id ON audit_logs(request_id);

-- Partition by month for better performance (optional)
-- CREATE INDEX idx_audit_logs_created_month ON audit_logs(date_trunc('month', created_date));

-- ============================================================
-- API_REQUEST_LOGS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS api_request_logs (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER,
    request_id          VARCHAR(255),
    method              VARCHAR(10) NOT NULL,
    path                VARCHAR(1000) NOT NULL,
    query_params        JSONB,
    request_headers     JSONB,
    request_body        JSONB,
    response_status     INTEGER,
    response_body       JSONB,
    ip_address          VARCHAR(45),
    user_agent          TEXT,
    duration_ms         INTEGER,
    error_message       TEXT,
    created_date        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_api_request_logs_user_id ON api_request_logs(user_id);
CREATE INDEX idx_api_request_logs_request_id ON api_request_logs(request_id);
CREATE INDEX idx_api_request_logs_path ON api_request_logs(path);
CREATE INDEX idx_api_request_logs_response_status ON api_request_logs(response_status);
CREATE INDEX idx_api_request_logs_created_date ON api_request_logs(created_date);

-- ============================================================
-- ERROR_LOGS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS error_logs (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER,
    request_id          VARCHAR(255),
    error_type          VARCHAR(100),
    error_code          VARCHAR(50),
    error_message       TEXT NOT NULL,
    stack_trace         TEXT,
    context             JSONB,
    path                VARCHAR(1000),
    ip_address          VARCHAR(45),
    user_agent          TEXT,
    is_resolved         BOOLEAN DEFAULT FALSE,
    resolved_at         TIMESTAMPTZ,
    resolved_by         INTEGER REFERENCES users(id),
    resolution_notes    TEXT,
    created_date        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_error_logs_user_id ON error_logs(user_id);
CREATE INDEX idx_error_logs_error_type ON error_logs(error_type);
CREATE INDEX idx_error_logs_error_code ON error_logs(error_code);
CREATE INDEX idx_error_logs_is_resolved ON error_logs(is_resolved);
CREATE INDEX idx_error_logs_created_date ON error_logs(created_date);

-- ============================================================
-- SCHEDULED_JOBS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS scheduled_jobs (
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR(100) NOT NULL,
    code                VARCHAR(100) UNIQUE NOT NULL,
    description         TEXT,
    cron_expression     VARCHAR(100),
    handler             VARCHAR(255) NOT NULL,
    parameters          JSONB DEFAULT '{}',
    is_active           BOOLEAN DEFAULT TRUE,
    last_run_at         TIMESTAMPTZ,
    last_run_status     VARCHAR(50),
    last_run_duration   INTEGER,
    last_error          TEXT,
    next_run_at         TIMESTAMPTZ,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_date        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_scheduled_jobs_code ON scheduled_jobs(code);
CREATE INDEX idx_scheduled_jobs_is_active ON scheduled_jobs(is_active);
CREATE INDEX idx_scheduled_jobs_next_run ON scheduled_jobs(next_run_at);

SELECT create_audit_triggers('scheduled_jobs');

-- ============================================================
-- JOB_EXECUTIONS TABLE
-- ============================================================

CREATE TABLE IF NOT EXISTS job_executions (
    id                  SERIAL PRIMARY KEY,
    job_id              INTEGER NOT NULL REFERENCES scheduled_jobs(id) ON DELETE CASCADE,
    status              VARCHAR(50) NOT NULL,
    started_at          TIMESTAMPTZ NOT NULL,
    finished_at         TIMESTAMPTZ,
    duration_ms         INTEGER,
    result              JSONB,
    error_message       TEXT,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_job_status CHECK (status IN ('running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX idx_job_executions_job_id ON job_executions(job_id);
CREATE INDEX idx_job_executions_status ON job_executions(status);
CREATE INDEX idx_job_executions_started_at ON job_executions(started_at);
