package service

import (
	"errors"
	"fmt"

	"github.com/go4s/iam/internal/model"
	"github.com/go4s/iam/internal/repository"
)

type PermissionService struct {
	permRepo           *repository.PermissionRepository
	rolePermissionRepo *repository.RolePermissionRepository
	roleRepo           *repository.RoleRepository
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		permRepo:           &repository.PermissionRepository{},
		rolePermissionRepo: &repository.RolePermissionRepository{},
		roleRepo:           &repository.RoleRepository{},
	}
}

// ListPermissions 列出权限（摘要格式）
func (s *PermissionService) ListPermissions(page, size int) ([]map[string]any, int64, error) {
	perms, total, err := s.permRepo.List(page, size)
	if err != nil {
		return nil, 0, err
	}

	fields, ok := GetFormat("permission", "summary")
	if !ok {
		return nil, 0, errors.New("permission summary format not found")
	}

	var result []map[string]any
	for _, p := range perms {
		data := map[string]any{
			"id":   fmt.Sprintf("permission:%d", p.ID),
			"name": p.Name,
			"code": p.Code,
		}
		result = append(result, FormatEntity("permission", fields, data))
	}
	return result, total, nil
}

// FocusPermission 聚焦权限（详情格式）
func (s *PermissionService) FocusPermission(id int64) (map[string]any, error) {
	perm, err := s.permRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, errors.New("permission not found")
	}

	fields, ok := GetFormat("permission", "detail")
	if !ok {
		return nil, errors.New("permission detail format not found")
	}

	// 获取拥有此权限的角色
	roleIDs, err := s.rolePermissionRepo.GetRoleIDsByPermissionID(perm.ID)
	if err != nil {
		return nil, err
	}

	var roles []string
	for _, rid := range roleIDs {
		role, err := s.roleRepo.GetByID(rid)
		if err != nil || role == nil {
			continue
		}
		roles = append(roles, fmt.Sprintf("role:%s", role.Code))
	}

	data := map[string]any{
		"id":          fmt.Sprintf("permission:%d", perm.ID),
		"name":        perm.Name,
		"code":        perm.Code,
		"resource":    perm.Resource,
		"action":      perm.Action,
		"description": perm.Description,
		"roles:":      roles,
	}

	result := FormatEntity("permission", fields, data)
	result["commands"] = s.buildPermissionCommands(perm.ID)
	return result, nil
}

// buildPermissionCommands 构建权限命令列表
func (s *PermissionService) buildPermissionCommands(permID int64) []map[string]any {
	return []map[string]any{
		{
			"action": "create_permission",
			"label":  "创建权限",
			"params": []map[string]any{
				{"name": "name", "type": "string", "required": true, "label": "权限名称"},
				{"name": "code", "type": "string", "required": true, "label": "权限编码"},
				{"name": "resource", "type": "string", "required": true, "label": "资源"},
				{"name": "action", "type": "string", "required": true, "label": "操作"},
				{"name": "description", "type": "string", "required": false, "label": "描述"},
			},
			"trigger": fmt.Sprintf("/api/v1/permission/%d/commands", permID),
			"sync":    true,
		},
	}
}

// CreatePermission 创建新权限
func (s *PermissionService) CreatePermission(callerID int64, params map[string]any) (map[string]any, error) {
	name, _ := params["name"].(string)
	code, _ := params["code"].(string)
	resource, _ := params["resource"].(string)
	action, _ := params["action"].(string)

	if name == "" || code == "" || resource == "" || action == "" {
		return nil, errors.New("name, code, resource, action are required")
	}

	// 检查权限码是否已存在
	existing, _ := s.permRepo.GetByCode(code)
	if existing != nil {
		return nil, errors.New("permission code already exists")
	}

	description, _ := params["description"].(string)
	perm := &model.Permission{
		Name:        name,
		Code:        code,
		Resource:    resource,
		Action:      action,
		Description: description,
	}
	if err := s.permRepo.Create(perm); err != nil {
		return nil, err
	}

	return map[string]any{
		"id":      fmt.Sprintf("permission:%d", perm.ID),
		"name":    perm.Name,
		"code":    perm.Code,
		"message": "权限创建成功",
	}, nil
}

// ExecutePermissionCommand 执行权限命令
func (s *PermissionService) ExecutePermissionCommand(callerID int64, targetID int64, action string, params map[string]any) (map[string]any, error) {
	// 检查权限
	requiredPerm := "permission:manage"
	userRoleRepo := &repository.UserRoleRepository{}
	rolePermRepo := &repository.RolePermissionRepository{}
	roleIDs, err := userRoleRepo.GetRoleIDsByUserID(callerID)
	if err != nil {
		return nil, err
	}

	hasPerm := false
	for _, roleID := range roleIDs {
		permIDs, _ := rolePermRepo.GetPermissionIDsByRoleID(roleID)
		for _, permID := range permIDs {
			perm, _ := s.permRepo.GetByID(permID)
			if perm != nil && perm.Code == requiredPerm {
				hasPerm = true
				break
			}
		}
		if hasPerm {
			break
		}
	}
	if !hasPerm {
		return nil, errors.New("permission denied: " + requiredPerm)
	}

	switch action {
	case "create_permission":
		return s.CreatePermission(callerID, params)
	default:
		return nil, errors.New("unknown command action")
	}
}
