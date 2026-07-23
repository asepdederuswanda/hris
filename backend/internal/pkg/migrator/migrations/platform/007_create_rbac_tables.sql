-- Migration: 007_create_rbac_tables
-- Database: Platform (Cross-Dialect)
-- Tabel untuk RBAC berbasis database: roles, permissions, role_permissions

CREATE TABLE IF NOT EXISTS rbac_roles (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    slug        VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255) NULL,
    parent_id   CHAR(36) NULL,
    is_system   SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_rbac_role_parent FOREIGN KEY (parent_id) REFERENCES rbac_roles(id) ON DELETE SET NULL
);

CREATE INDEX idx_rbac_roles_slug ON rbac_roles (slug);
CREATE INDEX idx_rbac_roles_parent ON rbac_roles (parent_id);
CREATE INDEX idx_rbac_roles_is_system ON rbac_roles (is_system);

CREATE TABLE IF NOT EXISTS rbac_permissions (
    id          CHAR(36) PRIMARY KEY,
    resource    VARCHAR(100) NOT NULL,
    action      VARCHAR(50) NOT NULL,
    description VARCHAR(255) NULL,
    is_system   SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uq_rbac_permission (resource, action)
);

CREATE INDEX idx_rbac_permissions_resource ON rbac_permissions (resource);
CREATE INDEX idx_rbac_permissions_is_system ON rbac_permissions (is_system);

CREATE TABLE IF NOT EXISTS rbac_role_permissions (
    role_id       CHAR(36) NOT NULL,
    permission_id CHAR(36) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (role_id, permission_id),

    CONSTRAINT fk_rbac_rp_role       FOREIGN KEY (role_id)       REFERENCES rbac_roles(id)       ON DELETE CASCADE,
    CONSTRAINT fk_rbac_rp_permission FOREIGN KEY (permission_id) REFERENCES rbac_permissions(id) ON DELETE CASCADE
);

CREATE INDEX idx_rbac_rp_role ON rbac_role_permissions (role_id);
CREATE INDEX idx_rbac_rp_permission ON rbac_role_permissions (permission_id);
