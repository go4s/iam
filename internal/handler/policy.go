package handler

import (
    "net/http"
    
    "github.com/go4s/iam/internal/auth"
    
    "github.com/gin-gonic/gin"
)

type PolicyHandler struct{}

type PolicyRequest struct {
    Sub string `json:"sub" binding:"required"`
    Obj string `json:"obj" binding:"required"`
    Act string `json:"act" binding:"required"`
}

func (h *PolicyHandler) AddPolicy(c *gin.Context) {
    var req PolicyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    added, err := auth.Enforcer.AddPolicy(req.Sub, req.Obj, req.Act)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if !added {
        c.JSON(http.StatusConflict, gin.H{"message": "Policy already exists"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Policy added successfully"})
}

func (h *PolicyHandler) RemovePolicy(c *gin.Context) {
    var req PolicyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    removed, err := auth.Enforcer.RemovePolicy(req.Sub, req.Obj, req.Act)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if !removed {
        c.JSON(http.StatusNotFound, gin.H{"message": "Policy not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Policy removed successfully"})
}

type GroupingRequest struct {
    Sub  string `json:"sub" binding:"required"`
    Role string `json:"role" binding:"required"`
}

func (h *PolicyHandler) AddGroupingPolicy(c *gin.Context) {
    var req GroupingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    added, err := auth.Enforcer.AddGroupingPolicy(req.Sub, req.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if !added {
        c.JSON(http.StatusConflict, gin.H{"message": "Grouping policy already exists"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Grouping policy added successfully"})
}

func (h *PolicyHandler) GetPolicies(c *gin.Context) {
    policies, err := auth.Enforcer.GetPolicy()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"policies": policies})
}

func (h *PolicyHandler) GetGroupingPolicies(c *gin.Context) {
    groupingPolicies, err := auth.Enforcer.GetGroupingPolicy()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"grouping_policies": groupingPolicies})
}

type EnforceRequest struct {
    Sub string `json:"sub" binding:"required"`
    Obj string `json:"obj" binding:"required"`
    Act string `json:"act" binding:"required"`
}

func (h *PolicyHandler) Enforce(c *gin.Context) {
    var req EnforceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    allowed, err := auth.Enforcer.Enforce(req.Sub, req.Obj, req.Act)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"allowed": allowed})
}
