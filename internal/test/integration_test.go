package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuth 测试认证接口
func TestAuth(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	t.Run("Login Success", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/auth/login", map[string]string{
			"username": "admin",
			"password": "admin123",
		}, "")

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		if _, ok := data["token"]; !ok {
			t.Error("Token not found")
		}
		if _, ok := data["user"]; !ok {
			t.Error("User not found")
		}
	})

	t.Run("Login Failure", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/auth/login", map[string]string{
			"username": "admin",
			"password": "wrongpassword",
		}, "")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Get Me", func(t *testing.T) {
		token := LoginAndGetToken(t, router, "admin", "admin123")
		w := MakeRequest(t, router, "GET", "/api/v1/auth/me", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		if data["username"] != "admin" {
			t.Errorf("Expected username admin, got %v", data["username"])
		}
	})

	t.Run("Get Me Without Token", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/auth/me", nil, "")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

// TestUser 测试用户接口
func TestUser(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("List Users", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user?page=1&size=10", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		items, ok := data["items"].([]any)
		if !ok {
			t.Fatal("Items not found")
		}

		if len(items) != 2 {
			t.Errorf("Expected 2 users, got %d", len(items))
		}
	})

	t.Run("Focus User", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user/1", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		if data["id"] != "user:1" {
			t.Errorf("Expected id user:1, got %v", data["id"])
		}

		// 检查 commands 是否存在
		if _, ok := data["commands"]; !ok {
			t.Error("Commands not found")
		}
	})

	t.Run("Focus Non-existent User", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user/999", nil, token)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("Command: Create User", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/user/1/commands", map[string]any{
			"action": "create_user",
			"params": map[string]any{
				"username": "testuser",
				"password": "testpass",
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})

	t.Run("Command: Reset Password", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/user/1/commands", map[string]any{
			"action": "reset_password",
			"params": map[string]any{
				"new_password": "newpass123",
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})

		t.Run("Command: Disable User", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/user/2/commands", map[string]any{
			"action": "disable",
			"params": map[string]any{},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 验证用户已禁用（无法登录）
		w = MakeRequest(t, router, "POST", "/api/v1/auth/login", map[string]string{
			"username": "editor",
			"password": "admin123",
		}, "")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for disabled user, got %d", w.Code)
		}
	})

	t.Run("Command: Self-Disable Forbidden", func(t *testing.T) {
		// admin 不能禁用自己
		w := MakeRequest(t, router, "POST", "/api/v1/user/1/commands", map[string]any{
			"action": "disable",
			"params": map[string]any{},
		}, token)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for self-disable, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertError(t, resp, "1004")
	})
}

// TestRole 测试角色接口
func TestRole(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("List Roles", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/role?page=1&size=10", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		items, ok := data["items"].([]any)
		if !ok {
			t.Fatal("Items not found")
		}

		if len(items) != 2 {
			t.Errorf("Expected 2 roles, got %d", len(items))
		}
	})

	t.Run("Focus Role", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/role/1", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		if data["id"] != "role:admin" {
			t.Errorf("Expected id role:admin, got %v", data["id"])
		}
	})

	t.Run("Command: Clone Role", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/role/1/commands", map[string]any{
			"action": "clone",
			"params": map[string]any{
				"new_name": "管理员副本",
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})
}

// TestPermission 测试权限接口
func TestPermission(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("List Permissions", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/permission?page=1&size=10", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		items, ok := data["items"].([]any)
		if !ok {
			t.Fatal("Items not found")
		}

		if len(items) != 4 {
			t.Errorf("Expected 4 permissions, got %d", len(items))
		}
	})

	t.Run("Focus Permission", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/permission/1", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})
}

// TestSystem 测试系统接口
func TestSystem(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("Reload Formats", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/system/reload-formats", nil, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})
}

// TestResponseFormat 测试统一响应格式
func TestResponseFormat(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	t.Run("Success Response Format", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/auth/login", map[string]string{
			"username": "admin",
			"password": "admin123",
		}, "")

		resp := ParseResponse(t, w)

		// 检查必须包含 code, message, data
		if _, ok := resp["code"]; !ok {
			t.Error("Response missing 'code'")
		}
		if _, ok := resp["message"]; !ok {
			t.Error("Response missing 'message'")
		}
		if _, ok := resp["data"]; !ok {
			t.Error("Response missing 'data'")
		}
	})

	t.Run("Error Response Format", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/auth/me", nil, "")

		resp := ParseResponse(t, w)

		if _, ok := resp["code"]; !ok {
			t.Error("Error response missing 'code'")
		}
		if _, ok := resp["message"]; !ok {
			t.Error("Error response missing 'message'")
		}
	})
}

// TestPagination 测试分页
func TestPagination(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("Pagination Format", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user?page=1&size=5", nil, token)

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		// 检查分页字段
		pagination, ok := data["pagination"].(map[string]any)
		if !ok {
			t.Fatal("Pagination not found")
		}

		if _, ok := pagination["page"]; !ok {
			t.Error("Pagination missing 'page'")
		}
		if _, ok := pagination["size"]; !ok {
			t.Error("Pagination missing 'size'")
		}
		if _, ok := pagination["total"]; !ok {
			t.Error("Pagination missing 'total'")
		}
	})
}

// TestEntityIDFormat 测试实体ID格式
func TestEntityIDFormat(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("User ID Format", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user/1", nil, token)

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		id, ok := data["id"].(string)
		if !ok || id != "user:1" {
			t.Errorf("Expected id 'user:1', got %v", data["id"])
		}
	})

	t.Run("Role ID Format", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/role/1", nil, token)

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		id, ok := data["id"].(string)
		if !ok || id != "role:admin" {
			t.Errorf("Expected id 'role:admin', got %v", data["id"])
		}
	})
}

// TestCommands 测试命令接口
func TestCommands(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("User Commands Structure", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user/1", nil, token)

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		commands, ok := data["commands"].([]any)
		if !ok {
			t.Fatal("Commands not found")
		}

		if len(commands) == 0 {
			t.Error("Expected commands, got empty")
		}

		// 检查命令结构
		for _, cmd := range commands {
			command, ok := cmd.(map[string]any)
			if !ok {
				t.Error("Invalid command format")
				continue
			}

			requiredFields := []string{"action", "label", "params", "trigger", "sync"}
			for _, field := range requiredFields {
				if _, ok := command[field]; !ok {
					t.Errorf("Command missing '%s'", field)
				}
			}
		}
	})

	t.Run("Permission Commands Empty", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/permission/1", nil, token)

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, ok := resp["data"].(map[string]any)
		if !ok {
			t.Fatal("Invalid data format")
		}

		commands, ok := data["commands"].([]any)
		if !ok {
			t.Fatal("Commands not found")
		}

		if len(commands) != 0 {
			t.Errorf("Expected empty commands, got %d", len(commands))
		}
	})
}

// TestRegression 回归测试：确保基本功能正常
func TestRegression(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	t.Run("Full Flow", func(t *testing.T) {
		// 1. 登录
		token := LoginAndGetToken(t, router, "admin", "admin123")
		fmt.Printf("Token obtained: %s...\n", token[:20])

		// 2. 获取当前用户
		w := MakeRequest(t, router, "GET", "/api/v1/auth/me", nil, token)
		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 3. 列出用户
		w = MakeRequest(t, router, "GET", "/api/v1/user", nil, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 4. 列出角色
		w = MakeRequest(t, router, "GET", "/api/v1/role", nil, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 5. 列出权限
		w = MakeRequest(t, router, "GET", "/api/v1/permission", nil, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 6. 创建用户
		w = MakeRequest(t, router, "POST", "/api/v1/user/1/commands", map[string]any{
			"action": "create_user",
			"params": map[string]any{
				"username": "regression_test",
				"password": "test123",
			},
		}, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 7. 验证新用户可登录
		newToken := LoginAndGetToken(t, router, "regression_test", "test123")
		if newToken == "" {
			t.Error("New user login failed")
		}

		fmt.Println("Regression test passed!")
	})
}

// TestBoundary 边界值测试
func TestBoundary(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("Page Size Zero", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user?page=1&size=0", nil, token)
		resp := ParseResponse(t, w)
		// size=0 时应该使用默认值
		AssertSuccess(t, resp)
	})

	t.Run("Page Size Too Large", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user?page=1&size=999", nil, token)
		resp := ParseResponse(t, w)
		// size 应该被限制在最大值
		AssertSuccess(t, resp)
	})

	t.Run("Negative Page", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user?page=-1&size=10", nil, token)
		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})

	t.Run("Invalid Entity ID", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user/abc", nil, token)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Empty Username", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/user/1/commands", map[string]any{
			"action": "create_user",
			"params": map[string]any{
				"username": "",
				"password": "password123",
			},
		}, token)
		resp := ParseResponse(t, w)
		AssertError(t, resp, "1002")
	})

	t.Run("Long Username", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/user/1/commands", map[string]any{
			"action": "create_user",
			"params": map[string]any{
				"username": "a very long username that exceeds normal limits and should still be handled gracefully by the system",
				"password": "password123",
			},
		}, token)
		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})

	t.Run("Empty Request Body", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/user/1/commands", nil, token)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/user/1/commands", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

// TestBusinessCommands 业务命令完整流程测试
func TestBusinessCommands(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	token := LoginAndGetToken(t, router, "admin", "admin123")

	t.Run("Assign Role to User", func(t *testing.T) {
		// 给用户2（editor）分配admin角色
		w := MakeRequest(t, router, "POST", "/api/v1/user/2/commands", map[string]any{
			"action": "assign_role",
			"params": map[string]any{
				"role_id": 1,
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 验证用户角色已更新
		w = MakeRequest(t, router, "GET", "/api/v1/user/2", nil, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, _ := resp["data"].(map[string]any)
		roles, _ := data["roles:"].([]any)

		foundAdmin := false
		for _, r := range roles {
			if r.(string) == "role:admin" {
				foundAdmin = true
				break
			}
		}

		if !foundAdmin {
			t.Error("Role assignment failed: admin role not found")
		}
	})

	t.Run("Add Permission to Role", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/role/2/commands", map[string]any{
			"action": "add_permission",
			"params": map[string]any{
				"permission_id": 1,
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 验证权限已添加
		w = MakeRequest(t, router, "GET", "/api/v1/role/2", nil, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, _ := resp["data"].(map[string]any)
		permissions, _ := data["permissions:"].([]any)

		if len(permissions) == 0 {
			t.Error("Permission add failed: no permissions found")
		}
	})

	t.Run("Remove Permission from Role", func(t *testing.T) {
		// 先添加权限
		MakeRequest(t, router, "POST", "/api/v1/role/2/commands", map[string]any{
			"action": "add_permission",
			"params": map[string]any{
				"permission_id": 2,
			},
		}, token)

		// 再移除权限
		w := MakeRequest(t, router, "POST", "/api/v1/role/2/commands", map[string]any{
			"action": "remove_permission",
			"params": map[string]any{
				"permission_id": 2,
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})

	t.Run("Clone Role and Verify", func(t *testing.T) {
		w := MakeRequest(t, router, "POST", "/api/v1/role/1/commands", map[string]any{
			"action": "clone",
			"params": map[string]any{
				"new_name": "克隆角色",
			},
		}, token)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 验证新角色存在
		w = MakeRequest(t, router, "GET", "/api/v1/role", nil, token)
		resp = ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, _ := resp["data"].(map[string]any)
		items, _ := data["items"].([]any)

		if len(items) != 3 {
			t.Errorf("Expected 3 roles after clone, got %d", len(items))
		}
	})

	t.Run("Reset Password and Login", func(t *testing.T) {
		// 重置密码
		w := MakeRequest(t, router, "POST", "/api/v1/user/2/commands", map[string]any{
			"action": "reset_password",
			"params": map[string]any{
				"new_password": "new_editor_pass",
			},
		}, token)

		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		// 使用新密码登录
		newToken := LoginAndGetToken(t, router, "editor", "new_editor_pass")
		if newToken == "" {
			t.Error("Login with new password failed")
		}
	})
}

// TestPermissionControl 权限控制测试
func TestPermissionControl(t *testing.T) {
	router := SetupTestServer()
	defer Teardown()

	adminToken := LoginAndGetToken(t, router, "admin", "admin123")
	editorToken := LoginAndGetToken(t, router, "editor", "admin123")

	t.Run("Editor Can Access User List", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user", nil, editorToken)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}
		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)
	})

	t.Run("Editor Cannot Execute Create User", func(t *testing.T) {
		// editor 尝试执行 create_user 命令应该被拒绝
		w := MakeRequest(t, router, "POST", "/api/v1/user/2/commands", map[string]any{
			"action": "create_user",
			"params": map[string]any{
				"username": "hacker",
				"password": "password123",
			},
		}, editorToken)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
		}

		resp := ParseResponse(t, w)
		AssertError(t, resp, "1004")
	})

	t.Run("Admin Has All Commands", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user/user:1", nil, adminToken)
		resp := ParseResponse(t, w)
		AssertSuccess(t, resp)

		data, _ := resp["data"].(map[string]any)
		commands, _ := data["commands"].([]any)

		requiredCommands := []string{"reset_password", "assign_role", "disable", "create_user"}
		found := make(map[string]bool)

		for _, cmd := range commands {
			command, _ := cmd.(map[string]any)
			action, _ := command["action"].(string)
			found[action] = true
		}

		for _, required := range requiredCommands {
			if !found[required] {
				t.Errorf("Admin missing command: %s", required)
			}
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		w := MakeRequest(t, router, "GET", "/api/v1/user", nil, "invalid_token")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		// 使用一个明显伪造的token
		w := MakeRequest(t, router, "GET", "/api/v1/user", nil, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/user", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}
