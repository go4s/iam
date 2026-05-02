package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{userService: service.NewUserService()}
}

// parseEntityID 从实体ID字符串解析数字ID，如 "user:1" -> 1
func parseEntityID(idStr string) (int64, error) {
	parts := strings.Split(idStr, ":")
	if len(parts) == 2 {
		return strconv.ParseInt(parts[1], 10, 64)
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func (h *UserHandler) List(c *gin.Context) {
	page, size := GetPageSize(c)
	keyword := GetKeyword(c)

	items, total, err := h.userService.ListUsers(page, size, keyword)
	if err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, err.Error(), nil)
		return
	}

	JSONList(c, items, page, size, total)
}

func (h *UserHandler) Focus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseEntityID(idStr)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid user id", nil)
		return
	}

	data, err := h.userService.FocusUser(id)
	if err != nil {
		JSONErr(c, http.StatusNotFound, CodeEntityNotFound, err.Error(), map[string]any{
			"template": "user",
			"id":       idStr,
		})
		return
	}

	JSONOK(c, data)
}

type CommandRequest struct {
	Action string         `json:"action" binding:"required"`
	Params map[string]any `json:"params"`
}

func (h *UserHandler) Command(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseEntityID(idStr)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, "invalid user id", nil)
		return
	}

	var req CommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	result, err := h.userService.ExecuteUserCommand(id, req.Action, req.Params)
	if err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	JSONOK(c, result)
}
