package middleware

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/go4s/iam/internal/auth"
)

func CasbinMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        username, exists := c.Get("username")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found"})
            c.Abort()
            return
        }
        
        obj := c.Request.URL.Path
        act := c.Request.Method
        
        allowed, err := auth.Enforcer.Enforce(username.(string), obj, act)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization check failed"})
            c.Abort()
            return
        }
        
        if !allowed {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
