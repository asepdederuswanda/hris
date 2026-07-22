-- =============================================================================
-- Tenant Migration: 010_approval
-- =============================================================================
-- Tabel untuk modul approval (workflow persetujuan) tenant.

-- ---------------------------------------------------------------------------
-- 10.1 Approval Flows (Master alur persetujuan)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_flows (
    id          CHAR(36) PRIMARY KEY,
    module      VARCHAR(255) NOT NULL COMMENT 'Modul tujuan: leave, overtime, adjustment, dll',
    name        VARCHAR(255) NOT NULL,
    version     INT NOT NULL DEFAULT 1,
    is_active   TINYINT(1) NOT NULL DEFAULT 1,
    deleted_at  TIMESTAMP NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_approval_flow_module (module),
    INDEX idx_approval_flow_active (is_active),
    INDEX idx_approval_flow_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 10.2 Approval Flow Steps (Langkah-langkah dalam alur)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_flow_steps (
    id                 CHAR(36) PRIMARY KEY,
    flow_id            CHAR(36) NOT NULL,
    step_order         INT NOT NULL,
    step_name          VARCHAR(100) NOT NULL,
    approver_type      ENUM('SUPERVISOR', 'ROLE', 'USER') NOT NULL,
    role_id            CHAR(36) NULL,
    approver_user_id   CHAR(36) NULL,
    approval_mode      ENUM('ANY_ONE', 'ALL', 'N_OF_M') NOT NULL DEFAULT 'ANY_ONE',
    required_approvals INT NULL,
    allow_reject       TINYINT(1) NOT NULL DEFAULT 1,
    conditions_json    JSON NULL,
    sla_hours          INT NULL,
    deleted_at         TIMESTAMP NULL,
    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_approval_step_order (flow_id, step_order),
    INDEX idx_approval_step_flow (flow_id),
    INDEX idx_approval_step_role (role_id),
    INDEX idx_approval_step_user (approver_user_id),
    INDEX idx_approval_step_deleted_at (deleted_at),

    CONSTRAINT fk_approval_step_flow FOREIGN KEY (flow_id) REFERENCES approval_flows(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 10.3 Approval Instances (Instance persetujuan per dokumen)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_instances (
    id            CHAR(36) PRIMARY KEY,
    module        VARCHAR(50) NOT NULL,
    document_id   CHAR(36) NOT NULL,
    flow_id       CHAR(36) NOT NULL,
    status        ENUM('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED') NOT NULL DEFAULT 'PENDING',
    current_step  INT NOT NULL DEFAULT 1,
    created_by    CHAR(36) NULL,
    deleted_at    TIMESTAMP NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_approval_instance_doc (module, document_id, deleted_at),
    INDEX idx_approval_instance_status (module, status),
    INDEX idx_approval_instance_flow (flow_id),
    INDEX idx_approval_instance_creator (created_by),
    INDEX idx_approval_instance_deleted_at (deleted_at),

    CONSTRAINT fk_approval_instance_flow FOREIGN KEY (flow_id) REFERENCES approval_flows(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 10.4 Approval Actions (Aksi yang dilakukan approver)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_actions (
    id            CHAR(36) PRIMARY KEY,
    instance_id   CHAR(36) NOT NULL,
    step_order    INT NOT NULL,
    actor_user_id CHAR(36) NOT NULL,
    action        ENUM('APPROVE', 'REJECT', 'CANCEL') NOT NULL,
    note          VARCHAR(255) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_approval_action_instance (instance_id, step_order, created_at),
    INDEX idx_approval_action_actor (actor_user_id, created_at),

    CONSTRAINT fk_approval_action_instance FOREIGN KEY (instance_id) REFERENCES approval_instances(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ---------------------------------------------------------------------------
-- 10.5 Approval Tasks (Task approval per approver)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS approval_tasks (
    id            CHAR(36) PRIMARY KEY,
    instance_id   CHAR(36) NOT NULL,
    step_order    INT NOT NULL,
    assignee_type ENUM('USER', 'ROLE') NOT NULL,
    assignee_id   CHAR(36) NOT NULL,
    status        ENUM('PENDING', 'DONE', 'CANCELLED') NOT NULL DEFAULT 'PENDING',
    deleted_at    TIMESTAMP NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_approval_task_assignee (assignee_type, assignee_id, status),
    INDEX idx_approval_task_instance (instance_id, step_order, status),
    INDEX idx_approval_task_deleted_at (deleted_at),

    CONSTRAINT fk_approval_task_instance FOREIGN KEY (instance_id) REFERENCES approval_instances(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
