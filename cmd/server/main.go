package main

import (
    "log"
    
    "github.com/go4s/iam/internal/auth"
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
    
    if err := auth.InitCasbin(); err != nil {
        log.Fatalf("Failed to initialize Casbin: %v", err)
    }
    
    userRepo := &repository.UserRepository{}
    authService := service.NewAuthService(userRepo)
    authHandler := handler.NewAuthHandler(authService)
    policyHandler := &handler.PolicyHandler{}
    
    r := gin.Default()
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "UP"})
    })
    
    authGroup := r.Group("/auth")
    {
        authGroup.POST("/register", authHandler.Register)
        authGroup.POST("/login", authHandler.Login)
        authGroup.POST("/validate", authHandler.Validate)
    }
    
    api := r.Group("/api/v1")
    api.Use(middleware.JWTMiddleware())
    api.Use(middleware.CasbinMiddleware())
    {
        api.POST("/enforce", policyHandler.Enforce)
        policies := api.Group("/policies")
        {
            policies.GET("", policyHandler.GetPolicies)
            policies.POST("", policyHandler.AddPolicy)
            policies.DELETE("", policyHandler.RemovePolicy)
            policies.GET("/grouping", policyHandler.GetGroupingPolicies)
            policies.POST("/grouping", policyHandler.AddGroupingPolicy)
        }
    }
    
    if err := r.Run(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
