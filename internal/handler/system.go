package handler

import (
	"fmt"
	"net/http"

	"github.com/go4s/iam/internal/repository"
	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// --- 1. 外部服务：新权限注册 ---

type RegisterPermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
}

func (h *SystemHandler) RegisterPermission(c *gin.Context) {
	callerIDRaw, exists := c.Get("user_id")
	if !exists {
		JSONErr(c, http.StatusUnauthorized, CodeUnauthorized, "user context not found", nil)
		return
	}
	callerID := int64(callerIDRaw.(float64))

	var req RegisterPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	// 检查 caller 是否有 permission:manage 权限
	if !checkUserPermission(callerID, "permission:manage") {
		JSONErr(c, http.StatusForbidden, CodeForbidden, "permission denied: permission:manage", nil)
		return
	}

	permService := service.NewPermissionService()
	result, err := permService.CreatePermission(callerID, map[string]any{
		"name":        req.Name,
		"code":        req.Code,
		"resource":    req.Resource,
		"action":      req.Action,
		"description": req.Description,
	})
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	JSONOK(c, result)
}

// --- 2. 外部服务：权限校验 ---

// CheckPermission 检查用户是否拥有指定权限
// GET /api/v1/system/check-permission?user_id=user:1&permission_code=article:delete
func (h *SystemHandler) CheckPermission(c *gin.Context) {
	permissionCode := c.Query("permission_code")
	if permissionCode == "" {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "permission_code is required", nil)
		return
	}

	var userID int64
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		// 从查询参数获取
		id, err := parseEntityID(userIDStr)
		if err != nil {
			JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid user_id", nil)
			return
		}
		userID = id
	} else {
		// 从 Token 解析
		callerIDRaw, exists := c.Get("user_id")
		if !exists {
			JSONErr(c, http.StatusUnauthorized, CodeUnauthorized, "user context not found", nil)
			return
		}
		userID = int64(callerIDRaw.(float64))
	}

	// 检查权限
	hasPermission := checkUserPermission(userID, permissionCode)

	JSONOK(c, map[string]any{
		"user_id":        fmt.Sprintf("user:%d", userID),
		"permission":     permissionCode,
		"has_permission": hasPermission,
	})
}

// --- 3. 外部服务：用户身份+权限校验 ---

type ValidateRequest struct {
	Token               string   `json:"token" binding:"required"`
	RequiredPermissions []string `json:"required_permissions"`
}

type ValidateResponse struct {
	Valid       bool     `json:"valid"`
	User        any      `json:"user,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Allowed     []string `json:"allowed,omitempty"`
	Denied      []string `json:"denied,omitempty"`
}

func (h *SystemHandler) ValidateIdentityAndPermission(c *gin.Context) {
	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	// 解析 Token
	authService := service.NewAuthService(&repository.UserRepository{})
	claims, err := authService.ParseToken(req.Token)
	if err != nil {
		JSONErr(c, http.StatusUnauthorized, CodeUnauthorized, "invalid token", nil)
		return
	}

	username := claims["sub"].(string)
	userIDFloat := claims["user_id"].(float64)
	userID := int64(userIDFloat)

	// 获取用户信息
	userInfo, err := authService.Me(username)
	if err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, sanitizeError(err), nil)
		return
	}

	// 获取用户所有权限
	allPerms := getUserPermissions(userID)

	// 检查所需权限
	var allowed []string
	var denied []string
	for _, required := range req.RequiredPermissions {
		has := false
		for _, p := range allPerms {
			if p == required {
				has = true
				break
			}
		}
		if has {
			allowed = append(allowed, required)
		} else {
			denied = append(denied, required)
		}
	}

	JSONOK(c, map[string]any{
		"valid":       true,
		"user":        userInfo,
		"permissions": allPerms,
		"allowed":     allowed,
		"denied":      denied,
	})
}

func (h *SystemHandler) ReloadFormats(c *gin.Context) {
	if err := service.ReloadFormats(); err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, sanitizeError(err), nil)
		return
	}
	JSONOK(c, map[string]any{"message": "formats reloaded"})
}

// --- 辅助函数 ---

// checkUserPermission 检查用户是否拥有指定权限
func checkUserPermission(userID int64, permissionCode string) bool {
	userRoleRepo := &repository.UserRoleRepository{}
	rolePermRepo := &repository.RolePermissionRepository{}
	permRepo := &repository.PermissionRepository{}

	roleIDs, err := userRoleRepo.GetRoleIDsByUserID(userID)
	if err != nil {
		return false
	}

	for _, roleID := range roleIDs {
		permIDs, _ := rolePermRepo.GetPermissionIDsByRoleID(roleID)
		for _, permID := range permIDs {
			perm, _ := permRepo.GetByID(permID)
			if perm != nil && perm.Code == permissionCode {
				return true
			}
		}
	}
	return false
}

// getUserPermissions 获取用户的所有权限码
func getUserPermissions(userID int64) []string {
	userRoleRepo := &repository.UserRoleRepository{}
	rolePermRepo := &repository.RolePermissionRepository{}
	permRepo := &repository.PermissionRepository{}

	var permissions []string
	seen := make(map[string]bool)

	roleIDs, err := userRoleRepo.GetRoleIDsByUserID(userID)
	if err != nil {
		return permissions
	}

	for _, roleID := range roleIDs {
		permIDs, _ := rolePermRepo.GetPermissionIDsByRoleID(roleID)
		for _, permID := range permIDs {
			perm, _ := permRepo.GetByID(permID)
			if perm != nil && !seen[perm.Code] {
				seen[perm.Code] = true
				permissions = append(permissions, perm.Code)
			}
		}
	}
	return permissions
}
