package handler

import (
	"net/http"

	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) ReloadFormats(c *gin.Context) {
	if err := service.ReloadFormats(); err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, sanitizeError(err), nil)
		return
	}
	JSONOK(c, map[string]any{"message": "formats reloaded"})
}
