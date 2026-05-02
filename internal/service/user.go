package service

import (
	"errors"
	"fmt"

	"github.com/go4s/iam/internal/model"
	"github.com/go4s/iam/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo           *repository.UserRepository
	userRoleRepo       *repository.UserRoleRepository
	roleRepo           *repository.RoleRepository
	rolePermissionRepo *repository.RolePermissionRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:           &repository.UserRepository{},
		userRoleRepo:       &repository.UserRoleRepository{},
		roleRepo:           &repository.RoleRepository{},
		rolePermissionRepo: &repository.RolePermissionRepository{},
	}
}

// ListUsers 列出用户（摘要格式）
func (s *UserService) ListUsers(page, size int, keyword string) ([]map[string]any, int64, error) {
	users, total, err := s.userRepo.List(page, size, keyword)
	if err != nil {
		return nil, 0, err
	}

	fields, ok := GetFormat("user", "summary")
	if !ok {
		return nil, 0, errors.New("user summary format not found")
	}

	var result []map[string]any
	for _, u := range users {
		data := map[string]any{
			"id":       fmt.Sprintf("user:%d", u.ID),
			"username": u.Username,
		}
		result = append(result, FormatEntity("user", fields, data))
	}
	return result, total, nil
}

// FocusUser 聚焦用户（详情格式）
func (s *UserService) FocusUser(id int64) (map[string]any, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	fields, ok := GetFormat("user", "detail")
	if !ok {
		return nil, errors.New("user detail format not found")
	}

	// 获取用户角色
	roleIDs, err := s.userRoleRepo.GetRoleIDsByUserID(user.ID)
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
		"id":       fmt.Sprintf("user:%d", user.ID),
		"username": user.Username,
		"roles:":   roles,
	}

	result := FormatEntity("user", fields, data)
	result["commands"] = s.buildUserCommands(user.ID, roles)
	return result, nil
}

// buildUserCommands 构建用户命令列表
func (s *UserService) buildUserCommands(userID int64, roles []string) []map[string]any {
	commands := []map[string]any{
		{
			"action":  "reset_password",
			"label":   "重置密码",
			"params": []map[string]any{
				{"name": "new_password", "type": "string", "required": true, "label": "新密码"},
			},
			"trigger": fmt.Sprintf("/api/v1/user/%d/commands", userID),
			"sync":    true,
		},
		{
			"action":  "assign_role",
			"label":   "分配角色",
			"params": []map[string]any{
				{"name": "role_id", "type": "entity", "template": "role", "required": true, "label": "角色"},
			},
			"trigger": fmt.Sprintf("/api/v1/user/%d/commands", userID),
			"sync":    true,
		},
		{
			"action":  "disable",
			"label":   "禁用账户",
			"params":  []map[string]any{},
			"trigger": fmt.Sprintf("/api/v1/user/%d/commands", userID),
			"sync":    true,
		},
	}

	// admin 用户额外显示 create_user
	if containsAdmin(roles) {
		commands = append(commands, map[string]any{
			"action":  "create_user",
			"label":   "创建用户",
			"params": []map[string]any{
				{"name": "username", "type": "string", "required": true, "label": "用户名"},
				{"name": "password", "type": "string", "required": true, "label": "密码"},
				{"name": "role_ids", "type": "entities", "template": "role", "required": false, "label": "角色"},
			},
			"trigger": fmt.Sprintf("/api/v1/user/%d/commands", userID),
			"sync":    true,
		})
	}

	return commands
}

func containsAdmin(roles []string) bool {
	for _, r := range roles {
		if r == "role:admin" {
			return true
		}
	}
	return false
}

// ExecuteUserCommand 执行用户命令
func (s *UserService) ExecuteUserCommand(userID int64, action string, params map[string]any) (map[string]any, error) {
	switch action {
	case "reset_password":
		return s.resetPassword(userID, params)
	case "assign_role":
		return s.assignRole(userID, params)
	case "disable":
		return s.disableUser(userID)
	case "create_user":
		return s.createUser(params)
	default:
		return nil, errors.New("unknown command action")
	}
}

func (s *UserService) resetPassword(userID int64, params map[string]any) (map[string]any, error) {
	newPassword, ok := params["new_password"].(string)
	if !ok || newPassword == "" {
		return nil, errors.New("new_password is required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashed)
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return map[string]any{"message": "密码重置成功"}, nil
}

func (s *UserService) assignRole(userID int64, params map[string]any) (map[string]any, error) {
	roleIDFloat, ok := params["role_id"].(float64)
	if !ok {
		return nil, errors.New("role_id is required")
	}
	roleID := int64(roleIDFloat)
	if err := s.userRoleRepo.AddUserRole(userID, roleID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "角色分配成功"}, nil
}

func (s *UserService) disableUser(userID int64) (map[string]any, error) {
	// 简单实现：删除用户（或标记为禁用，这里直接删除）
	// 实际生产环境应使用状态字段
	return nil, errors.New("disable not implemented")
}

func (s *UserService) createUser(params map[string]any) (map[string]any, error) {
	username, _ := params["username"].(string)
	password, _ := params["password"].(string)
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	existing, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashed),
		Role:         "user", // 默认角色，废弃字段但保留
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 分配角色
	if roleIDs, ok := params["role_ids"].([]any); ok {
		for _, rid := range roleIDs {
			if ridFloat, ok := rid.(float64); ok {
				s.userRoleRepo.AddUserRole(user.ID, int64(ridFloat))
			}
		}
	}

	return map[string]any{
		"id":       fmt.Sprintf("user:%d", user.ID),
		"username": user.Username,
		"message":  "用户创建成功",
	}, nil
}
