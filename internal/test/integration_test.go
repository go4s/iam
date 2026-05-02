package test

import (
	"fmt"
	"net/http"
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
