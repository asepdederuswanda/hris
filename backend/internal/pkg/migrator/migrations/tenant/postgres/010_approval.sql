-- =============================================================================
-- Tenant Migration: 010_approval
-- =============================================================================
-- Tabel untuk modul approval (workflow persetujuan) tenant.

-- ---------------------------------------------------------------------------
-- 10.1 Approval Flows (Master alur persetujuan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_flows (
    id          CHAR(36) PRIMARY KEY,
    module      VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    version     INT NOT NULL DEFAULT 1,
    is_active   SMALLINT NOT NULL DEFAULT 1,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_approval_flow_module ON approval_flows (module);

CREATE INDEX IF NOT EXISTS idx_approval_flow_active ON approval_flows (is_active);

CREATE INDEX IF NOT EXISTS idx_approval_flow_deleted_at ON approval_flows (deleted_at);

-- ---------------------------------------------------------------------------
-- 10.2 Approval Flow Steps (Langkah-langkah dalam alur)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_flow_steps (
    id                 CHAR(36) PRIMARY KEY,
    flow_id            CHAR(36) NOT NULL,
    step_order         INT NOT NULL,
    step_name          VARCHAR(100) NOT NULL,
    approver_type      VARCHAR(255) NOT NULL,
    role_id            CHAR(36) NULL,
    approver_user_id   CHAR(36) NULL,
    approval_mode      VARCHAR(255) NOT NULL DEFAULT 'ANY_ONE',
    required_approvals INT NULL,
    allow_reject       SMALLINT NOT NULL DEFAULT 1,
    conditions_json    JSON NULL,
    sla_hours          INT NULL,
    deleted_at         TIMESTAMP,
    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_approval_step_order UNIQUE (flow_id, step_order),

    CONSTRAINT fk_approval_step_flow FOREIGN KEY (flow_id) REFERENCES approval_flows(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_approval_step_flow ON approval_flow_steps (flow_id);

CREATE INDEX IF NOT EXISTS idx_approval_step_role ON approval_flow_steps (role_id);

CREATE INDEX IF NOT EXISTS idx_approval_step_user ON approval_flow_steps (approver_user_id);

CREATE INDEX IF NOT EXISTS idx_approval_step_deleted_at ON approval_flow_steps (deleted_at);

-- ---------------------------------------------------------------------------
-- 10.3 Approval Instances (Instance persetujuan per dokumen)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_instances (
    id            CHAR(36) PRIMARY KEY,
    module        VARCHAR(50) NOT NULL,
    document_id   CHAR(36) NOT NULL,
    flow_id       CHAR(36) NOT NULL,
    status        VARCHAR(255) NOT NULL DEFAULT 'PENDING',
    current_step  INT NOT NULL DEFAULT 1,
    created_by    CHAR(36) NULL,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_approval_instance_doc UNIQUE (module, document_id, deleted_at),

    CONSTRAINT fk_approval_instance_flow FOREIGN KEY (flow_id) REFERENCES approval_flows(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_approval_instance_status ON approval_instances (module, status);

CREATE INDEX IF NOT EXISTS idx_approval_instance_flow ON approval_instances (flow_id);

CREATE INDEX IF NOT EXISTS idx_approval_instance_creator ON approval_instances (created_by);

CREATE INDEX IF NOT EXISTS idx_approval_instance_deleted_at ON approval_instances (deleted_at);

-- ---------------------------------------------------------------------------
-- 10.4 Approval Actions (Aksi yang dilakukan approver)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_actions (
    id            CHAR(36) PRIMARY KEY,
    instance_id   CHAR(36) NOT NULL,
    step_order    INT NOT NULL,
    actor_user_id CHAR(36) NOT NULL,
    action        VARCHAR(255) NOT NULL,
    note          VARCHAR(255) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_approval_action_instance FOREIGN KEY (instance_id) REFERENCES approval_instances(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_approval_action_instance ON approval_actions (instance_id, step_order, created_at);

CREATE INDEX IF NOT EXISTS idx_approval_action_actor ON approval_actions (actor_user_id, created_at);

-- ---------------------------------------------------------------------------
-- 10.5 Approval Tasks (Task approval per approver)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_tasks (
    id            CHAR(36) PRIMARY KEY,
    instance_id   CHAR(36) NOT NULL,
    step_order    INT NOT NULL,
    assignee_type VARCHAR(255) NOT NULL,
    assignee_id   CHAR(36) NOT NULL,
    status        VARCHAR(255) NOT NULL DEFAULT 'PENDING',
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_approval_task_instance FOREIGN KEY (instance_id) REFERENCES approval_instances(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_approval_task_assignee ON approval_tasks (assignee_type, assignee_id, status);

CREATE INDEX IF NOT EXISTS idx_approval_task_instance ON approval_tasks (instance_id, step_order, status);

CREATE INDEX IF NOT EXISTS idx_approval_task_deleted_at ON approval_tasks (deleted_at);
