package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 业务码常量
const (
	CodeSuccess          = "0000"
	CodeEntityNotFound   = "1001"
	CodeBadRequest       = "1002"
	CodeUnauthorized     = "1003"
	CodeForbidden        = "1004"
	CodeConflict         = "1005"
	CodeValidationFailed = "1006"
	CodeInternalError    = "9999"
)

// Response 统一响应结构
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Pagination 分页信息
type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// ListData 列表数据
type ListData struct {
	Items      any        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// JSONOK 返回成功响应
func JSONOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// JSONList 返回列表响应
func JSONList(c *gin.Context, items any, page, size int, total int64) {
	JSONOK(c, ListData{
		Items: items,
		Pagination: Pagination{
			Page:  page,
			Size:  size,
			Total: int(total),
		},
	})
}

// JSONErr 返回错误响应
func JSONErr(c *gin.Context, httpCode int, code string, message string, details any) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    details,
	})
}

// GetPageSize 从查询参数获取分页信息
func GetPageSize(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	size, _ = strconv.Atoi(c.DefaultQuery("size", "20"))
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}

// GetKeyword 获取搜索关键词
func GetKeyword(c *gin.Context) string {
	return c.DefaultQuery("q", "")
}
