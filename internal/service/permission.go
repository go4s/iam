package service

import (
	"errors"
	"fmt"

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
			"id":   fmt.Sprintf("permission:%s", p.Code),
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
		"id":          fmt.Sprintf("permission:%s", perm.Code),
		"name":        perm.Name,
		"code":        perm.Code,
		"resource":    perm.Resource,
		"action":      perm.Action,
		"description": perm.Description,
		"roles:":      roles,
	}

	result := FormatEntity("permission", fields, data)
	result["commands"] = []map[string]any{}
	return result, nil
}

// ExecutePermissionCommand 执行权限命令（空实现）
func (s *PermissionService) ExecutePermissionCommand(id int64, action string, params map[string]any) (map[string]any, error) {
	return nil, errors.New("permission does not support commands")
}
