package main

import (
	"log"

	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/handler"
	"github.com/go4s/iam/internal/middleware"
	"github.com/go4s/iam/internal/repository"
	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 加载实体格式配置
	if err := service.LoadFormats(); err != nil {
		log.Fatalf("Failed to load formats: %v", err)
	}

	userRepo := &repository.UserRepository{}
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	permHandler := handler.NewPermissionHandler()
	systemHandler := handler.NewSystemHandler()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	{
		api := r.Group("/api/v0")
		// 外部服务接入（无需额外认证，通过请求参数验证）
		api.POST("/system/validate", systemHandler.ValidateIdentityAndPermission)
	}
	// API v1 基础路径
	api := r.Group("/api/v1")
	api.Use(middleware.RateLimitMiddleware())

	// 认证接口（不需要 JWT，但会被限流）
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// 需要认证的接口
	authorized := api.Group("")
	authorized.Use(middleware.JWTMiddleware())
	{
		authorized.GET("/auth/me", authHandler.Me)

		// 用户接口
		authorized.GET("/user", userHandler.List)
		authorized.GET("/user/:id", userHandler.Focus)
		authorized.POST("/user/:id/commands", userHandler.Command)

		// 角色接口
		authorized.GET("/role", roleHandler.List)
		authorized.GET("/role/:id", roleHandler.Focus)
		authorized.POST("/role/:id/commands", roleHandler.Command)

		// 权限接口
		authorized.GET("/permission", permHandler.List)
		authorized.GET("/permission/:id", permHandler.Focus)
		authorized.POST("/permission/:id/commands", permHandler.Command)

		// 系统接口（需要认证）
		authorized.POST("/system/reload-formats", systemHandler.ReloadFormats)
		authorized.POST("/system/permissions", systemHandler.RegisterPermission)
		authorized.GET("/system/check-permission", systemHandler.CheckPermission)

	}

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
