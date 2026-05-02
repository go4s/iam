package service

import (
	"errors"
	"fmt"

	"github.com/go4s/iam/internal/model"
	"github.com/go4s/iam/internal/repository"
)

type RoleService struct {
	roleRepo           *repository.RoleRepository
	rolePermissionRepo *repository.RolePermissionRepository
	userRoleRepo       *repository.UserRoleRepository
	userRepo           *repository.UserRepository
	permRepo           *repository.PermissionRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo:           &repository.RoleRepository{},
		rolePermissionRepo: &repository.RolePermissionRepository{},
		userRoleRepo:       &repository.UserRoleRepository{},
		userRepo:           &repository.UserRepository{},
		permRepo:           &repository.PermissionRepository{},
	}
}

// ListRoles 列出角色（摘要格式）
func (s *RoleService) ListRoles(page, size int) ([]map[string]any, int64, error) {
	roles, total, err := s.roleRepo.List(page, size)
	if err != nil {
		return nil, 0, err
	}

	fields, ok := GetFormat("role", "summary")
	if !ok {
		return nil, 0, errors.New("role summary format not found")
	}

	var result []map[string]any
	for _, r := range roles {
		data := map[string]any{
			"id":   fmt.Sprintf("role:%s", r.Code),
			"name": r.Name,
			"code": r.Code,
		}
		result = append(result, FormatEntity("role", fields, data))
	}
	return result, total, nil
}

// FocusRole 聚焦角色（详情格式）
func (s *RoleService) FocusRole(id int64) (map[string]any, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	fields, ok := GetFormat("role", "detail")
	if !ok {
		return nil, errors.New("role detail format not found")
	}

	// 获取角色权限
	permIDs, err := s.rolePermissionRepo.GetPermissionIDsByRoleID(role.ID)
	if err != nil {
		return nil, err
	}

	var permissions []string
	for _, pid := range permIDs {
		perm, err := s.permRepo.GetByID(pid)
		if err != nil || perm == nil {
			continue
		}
		permissions = append(permissions, fmt.Sprintf("permission:%s", perm.Code))
	}

	// 获取角色下的用户
	userIDs, err := s.userRoleRepo.GetUserIDsByRoleID(role.ID)
	if err != nil {
		return nil, err
	}

	var users []string
	for _, uid := range userIDs {
		users = append(users, fmt.Sprintf("user:%d", uid))
	}

	data := map[string]any{
		"id":          fmt.Sprintf("role:%s", role.Code),
		"name":        role.Name,
		"code":        role.Code,
		"description": role.Description,
		"permissions:": permissions,
		"users:":       users,
	}

	result := FormatEntity("role", fields, data)
	result["commands"] = s.buildRoleCommands(role.ID)
	return result, nil
}

// buildRoleCommands 构建角色命令列表
func (s *RoleService) buildRoleCommands(roleID int64) []map[string]any {
	return []map[string]any{
		{
			"action":  "add_permission",
			"label":   "添加权限",
			"params": []map[string]any{
				{"name": "permission_id", "type": "entity", "template": "permission", "required": true, "label": "权限"},
			},
			"trigger": fmt.Sprintf("/api/v1/role/%d/commands", roleID),
			"sync":    true,
		},
		{
			"action":  "remove_permission",
			"label":   "移除权限",
			"params": []map[string]any{
				{"name": "permission_id", "type": "entity", "template": "permission", "required": true, "label": "权限"},
			},
			"trigger": fmt.Sprintf("/api/v1/role/%d/commands", roleID),
			"sync":    true,
		},
		{
			"action":  "clone",
			"label":   "复制角色",
			"params": []map[string]any{
				{"name": "new_name", "type": "string", "required": true, "label": "新角色名称"},
			},
			"trigger": fmt.Sprintf("/api/v1/role/%d/commands", roleID),
			"sync":    true,
		},
	}
}

// ExecuteRoleCommand 执行角色命令
func (s *RoleService) ExecuteRoleCommand(roleID int64, action string, params map[string]any) (map[string]any, error) {
	switch action {
	case "add_permission":
		return s.addPermission(roleID, params)
	case "remove_permission":
		return s.removePermission(roleID, params)
	case "clone":
		return s.cloneRole(roleID, params)
	default:
		return nil, errors.New("unknown command action")
	}
}

func (s *RoleService) addPermission(roleID int64, params map[string]any) (map[string]any, error) {
	permIDFloat, ok := params["permission_id"].(float64)
	if !ok {
		return nil, errors.New("permission_id is required")
	}
	permID := int64(permIDFloat)
	if err := s.rolePermissionRepo.AddRolePermission(roleID, permID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "权限添加成功"}, nil
}

func (s *RoleService) removePermission(roleID int64, params map[string]any) (map[string]any, error) {
	permIDFloat, ok := params["permission_id"].(float64)
	if !ok {
		return nil, errors.New("permission_id is required")
	}
	permID := int64(permIDFloat)
	if err := s.rolePermissionRepo.RemoveRolePermission(roleID, permID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "权限移除成功"}, nil
}

func (s *RoleService) cloneRole(roleID int64, params map[string]any) (map[string]any, error) {
	newName, ok := params["new_name"].(string)
	if !ok || newName == "" {
		return nil, errors.New("new_name is required")
	}

	srcRole, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if srcRole == nil {
		return nil, errors.New("source role not found")
	}

	newRole := &model.Role{
		Name:        newName,
		Code:        srcRole.Code + "_copy",
		Description: srcRole.Description,
	}
	if err := s.roleRepo.Create(newRole); err != nil {
		return nil, err
	}

	// 复制权限关联
	permIDs, _ := s.rolePermissionRepo.GetPermissionIDsByRoleID(roleID)
	for _, pid := range permIDs {
		s.rolePermissionRepo.AddRolePermission(newRole.ID, pid)
	}

	return map[string]any{
		"id":      fmt.Sprintf("role:%s", newRole.Code),
		"name":    newRole.Name,
		"message": "角色复制成功",
	}, nil
}
