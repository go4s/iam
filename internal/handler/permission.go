package handler

import (
	"net/http"

	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permService *service.PermissionService
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{permService: service.NewPermissionService()}
}

func (h *PermissionHandler) List(c *gin.Context) {
	page, size := GetPageSize(c)

	items, total, err := h.permService.ListPermissions(page, size)
	if err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, sanitizeError(err), nil)
		return
	}

	JSONList(c, items, page, size, total)
}

func (h *PermissionHandler) Focus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseEntityID(idStr)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid permission id", nil)
		return
	}

	data, err := h.permService.FocusPermission(id)
	if err != nil {
		JSONErr(c, http.StatusNotFound, CodeEntityNotFound, sanitizeError(err), map[string]any{
			"template": "permission",
			"id":       idStr,
		})
		return
	}

	JSONOK(c, data)
}

func (h *PermissionHandler) Command(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseEntityID(idStr)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid permission id", nil)
		return
	}

	var req CommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	result, err := h.permService.ExecutePermissionCommand(id, req.Action, req.Params)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, sanitizeError(err), nil)
		return
	}

	JSONOK(c, result)
}
