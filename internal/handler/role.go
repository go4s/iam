package handler

import (
	"net/http"

	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{roleService: service.NewRoleService()}
}

func (h *RoleHandler) List(c *gin.Context) {
	page, size := GetPageSize(c)

	items, total, err := h.roleService.ListRoles(page, size)
	if err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, err.Error(), nil)
		return
	}

	JSONList(c, items, page, size, total)
}

func (h *RoleHandler) Focus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseEntityID(idStr)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid role id", nil)
		return
	}

	data, err := h.roleService.FocusRole(id)
	if err != nil {
		JSONErr(c, http.StatusNotFound, CodeEntityNotFound, err.Error(), map[string]any{
			"template": "role",
			"id":       idStr,
		})
		return
	}

	JSONOK(c, data)
}

func (h *RoleHandler) Command(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseEntityID(idStr)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid role id", nil)
		return
	}

	var req CommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	result, err := h.roleService.ExecuteRoleCommand(id, req.Action, req.Params)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	JSONOK(c, result)
}
