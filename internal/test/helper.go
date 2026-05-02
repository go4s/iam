package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go4s/iam/internal/handler"
	"github.com/go4s/iam/internal/middleware"
	"github.com/go4s/iam/internal/repository"
	"github.com/go4s/iam/internal/service"
)

var (
	testRouter *gin.Engine
	testToken  string
)

// SetupTestServer 创建测试服务器
func SetupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// 初始化数据库
	if err := SetupTestDB(); err != nil {
		panic(err)
	}

	// 加载测试数据
	if err := LoadFixtures(); err != nil {
		panic(err)
	}

	// 加载格式配置
	if err := service.LoadFormats(); err != nil {
		panic(err)
	}

	// 创建 Handler
	userRepo := &repository.UserRepository{}
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	permHandler := handler.NewPermissionHandler()
	systemHandler := handler.NewSystemHandler()

	r := gin.New()

	// API v1
	api := r.Group("/api/v1")
	{
		// 认证接口
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
		}

		// 需要认证的接口
		authorized := api.Group("")
		authorized.Use(middleware.JWTMiddleware())
		{
			authorized.GET("/auth/me", authHandler.Me)
			authorized.GET("/user", userHandler.List)
			authorized.GET("/user/:id", userHandler.Focus)
			authorized.POST("/user/:id/commands", userHandler.Command)
			authorized.GET("/role", roleHandler.List)
			authorized.GET("/role/:id", roleHandler.Focus)
			authorized.POST("/role/:id/commands", roleHandler.Command)
			authorized.GET("/permission", permHandler.List)
			authorized.GET("/permission/:id", permHandler.Focus)
			authorized.POST("/permission/:id/commands", permHandler.Command)
			authorized.POST("/system/reload-formats", systemHandler.ReloadFormats)
		}
	}

	return r
}

// LoginAndGetToken 登录并获取 token
func LoginAndGetToken(t *testing.T, router *gin.Engine, username, password string) string {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed: %d %s", w.Code, w.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	data, ok := response["data"].(map[string]any)
	if !ok {
		t.Fatalf("Invalid response format: %v", response)
	}

	token, ok := data["token"].(string)
	if !ok {
		t.Fatalf("Token not found in response: %v", data)
	}

	return token
}

// MakeRequest 发送 HTTP 请求
func MakeRequest(t *testing.T, router *gin.Engine, method, path string, body any, token string) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	// 设置 RemoteAddr 避免限流中间件 panic
	req.RemoteAddr = "127.0.0.1:12345"

	router.ServeHTTP(w, req)
	return w
}

// ParseResponse 解析响应
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v\nBody: %s", err, w.Body.String())
	}
	return response
}

// AssertSuccess 断言成功响应
func AssertSuccess(t *testing.T, response map[string]any) {
	code, ok := response["code"].(string)
	if !ok || code != "0000" {
		t.Fatalf("Expected success code 0000, got: %v", response)
	}
}

// AssertError 断言错误响应
func AssertError(t *testing.T, response map[string]any, expectedCode string) {
	code, ok := response["code"].(string)
	if !ok || code != expectedCode {
		t.Fatalf("Expected error code %s, got: %v", expectedCode, response)
	}
}

// Teardown 清理测试环境
func Teardown() {
	CleanupTestDB()
}
