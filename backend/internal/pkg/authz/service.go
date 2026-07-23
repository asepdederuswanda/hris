package authz

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service untuk manajemen RBAC (roles & permissions).
type Service struct {
	repo     *Repository
	enforcer *Enforcer
	logger   *zap.Logger
}

// NewService membuat Service RBAC baru.
// Setelah perubahan role/permission, panggil service.Sync() untuk
// menyegarkan Enforcer tanpa restart server.
func NewService(repo *Repository, enforcer *Enforcer, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		enforcer: enforcer,
		logger:   logger,
	}
}

// Sync menyegarkan Enforcer dari database. Panggil setelah ada perubahan.
func (s *Service) Sync() error {
	if err := s.enforcer.Reload(); err != nil {
		return fmt.Errorf("failed to reload enforcer: %w", err)
	}
	s.logger.Info("RBAC enforcer reloaded from database")
	return nil
}

// =========================================================================
// Role Management
// =========================================================================

type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Slug        string  `json:"slug" binding:"required,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	ParentID    *string `json:"parent_id"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	ParentID    *string `json:"parent_id"`
}

type RoleResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description,omitempty"`
	ParentID    string           `json:"parent_id,omitempty"`
	IsSystem    bool             `json:"is_system"`
	Permissions []PermissionInfo `json:"permissions,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type PermissionInfo struct {
	ID          string `json:"id"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
}

func (s *Service) CreateRole(req CreateRoleRequest) (*RoleResponse, error) {
	var parentID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentID = &pid
	}

	role := &RbacRole{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    parentID,
	}

	if err := s.repo.CreateRole(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	resp := roleToResponse(role, nil)
	return &resp, nil
}

func (s *Service) ListRoles() ([]RoleResponse, error) {
	roles, err := s.repo.FindAllRoles()
	if err != nil {
		return nil, err
	}

	// Load all permissions for enrichment
	perms, _ := s.repo.FindAllPermissions()
	permMap := make(map[string]PermissionInfo)
	for _, p := range perms {
		permMap[p.ID.String()] = PermissionInfo{
			ID:          p.ID.String(),
			Resource:    p.Resource,
			Action:      p.Action,
			Description: safeStr(p.Description),
		}
	}

	var responses []RoleResponse
	for _, role := range roles {
		var rolePerms []PermissionInfo
		for _, rp := range role.Permissions {
			if pi, ok := permMap[rp.PermissionID.String()]; ok {
				rolePerms = append(rolePerms, pi)
			}
		}
		responses = append(responses, roleToResponse(&role, rolePerms))
	}
	return responses, nil
}

func (s *Service) GetRole(id string) (*RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid role id: %w", err)
	}

	role, err := s.repo.FindRoleByID(uid)
	if err != nil {
		return nil, err
	}

	// Load permissions for this role
	rps, _ := s.repo.FindRolePermissions(uid)
	var rolePerms []PermissionInfo
	for _, rp := range rps {
		perm, err := s.repo.FindPermissionByID(rp.PermissionID)
		if err != nil {
			continue
		}
		rolePerms = append(rolePerms, PermissionInfo{
			ID:          perm.ID.String(),
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: safeStr(perm.Description),
		})
	}

	resp := roleToResponse(role, rolePerms)
	return &resp, nil
}

func (s *Service) UpdateRole(id string, req UpdateRoleRequest) (*RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid role id: %w", err)
	}

	role, err := s.repo.FindRoleByID(uid)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = req.Description
	}
	if req.ParentID != nil {
		if *req.ParentID == "" {
			role.ParentID = nil
		} else {
			pid, err := uuid.Parse(*req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("invalid parent_id: %w", err)
			}
			role.ParentID = &pid
		}
	}

	if err := s.repo.UpdateRole(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	s.logger.Info("Role updated", zap.String("role", role.Slug))
	resp := roleToResponse(role, nil)
	return &resp, nil
}

func (s *Service) DeleteRole(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid role id: %w", err)
	}
	return s.repo.DeleteRole(uid)
}

// =========================================================================
// Permission Management
// =========================================================================

type CreatePermissionRequest struct {
	Resource    string  `json:"resource" binding:"required,max=100"`
	Action      string  `json:"action" binding:"required,max=50"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type PermissionResponse struct {
	ID          string    `json:"id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Service) CreatePermission(req CreatePermissionRequest) (*PermissionResponse, error) {
	perm := &RbacPermission{
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := s.repo.CreatePermission(perm); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	resp := PermissionResponse{
		ID:          perm.ID.String(),
		Resource:    perm.Resource,
		Action:      perm.Action,
		Description: safeStr(perm.Description),
		IsSystem:    perm.IsSystem,
		CreatedAt:   perm.CreatedAt,
	}
	return &resp, nil
}

func (s *Service) ListPermissions() ([]PermissionResponse, error) {
	perms, err := s.repo.FindAllPermissions()
	if err != nil {
		return nil, err
	}

	var responses []PermissionResponse
	for _, p := range perms {
		responses = append(responses, PermissionResponse{
			ID:          p.ID.String(),
			Resource:    p.Resource,
			Action:      p.Action,
			Description: safeStr(p.Description),
			IsSystem:    p.IsSystem,
			CreatedAt:   p.CreatedAt,
		})
	}
	return responses, nil
}

func (s *Service) DeletePermission(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid permission id: %w", err)
	}
	return s.repo.DeletePermission(uid)
}

// =========================================================================
// Role-Permission Assignment
// =========================================================================

type AssignPermissionRequest struct {
	PermissionID string `json:"permission_id" binding:"required"`
}

func (s *Service) AssignPermission(roleID string, permissionID string) error {
	rUID, err := uuid.Parse(roleID)
	if err != nil {
		return fmt.Errorf("invalid role id: %w", err)
	}
	pUID, err := uuid.Parse(permissionID)
	if err != nil {
		return fmt.Errorf("invalid permission id: %w", err)
	}

	if err := s.repo.AssignPermission(rUID, pUID); err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	s.logger.Info("Permission assigned to role",
		zap.String("role_id", roleID),
		zap.String("permission_id", permissionID))

	return nil
}

func (s *Service) RevokePermission(roleID string, permissionID string) error {
	rUID, err := uuid.Parse(roleID)
	if err != nil {
		return fmt.Errorf("invalid role id: %w", err)
	}
	pUID, err := uuid.Parse(permissionID)
	if err != nil {
		return fmt.Errorf("invalid permission id: %w", err)
	}

	if err := s.repo.RevokePermission(rUID, pUID); err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	s.logger.Info("Permission revoked from role",
		zap.String("role_id", roleID),
		zap.String("permission_id", permissionID))

	return nil
}

// =========================================================================
// Helpers
// =========================================================================

func roleToResponse(role *RbacRole, perms []PermissionInfo) RoleResponse {
	resp := RoleResponse{
		ID:        role.ID.String(),
		Name:      role.Name,
		Slug:      role.Slug,
		IsSystem:  role.IsSystem,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
	if role.Description != nil {
		resp.Description = *role.Description
	}
	if role.ParentID != nil {
		resp.ParentID = role.ParentID.String()
	}
	if perms != nil {
		resp.Permissions = perms
	}
	return resp
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
