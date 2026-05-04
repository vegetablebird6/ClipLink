package response

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/xiaojiu/cliplink/internal/common/apperr"
	"github.com/xiaojiu/cliplink/internal/common/i18n"
)

// 状态码常量
const (
	StatusSuccess      = http.StatusOK
	StatusBadRequest   = http.StatusBadRequest
	StatusUnauthorized = http.StatusUnauthorized
	StatusForbidden    = http.StatusForbidden
	StatusNotFound     = http.StatusNotFound
	StatusServerError  = http.StatusInternalServerError
)

// Response 统一响应结构
type Response struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Success    bool        `json:"success"`
	ErrorCode  string      `json:"error_code,omitempty"`
	MessageKey string      `json:"message_key,omitempty"`
	Details    string      `json:"details,omitempty"`
}

// normalizeSlice 将 nil、nil 指针、nil 切片统一处理：
// - nil 接口值 → 保持 nil（由 omitempty 处理）
// - nil 指针 → 保持 nil
// - nil 切片 → 转换为空切片 []interface{}{}，避免 JSON 输出 null
func normalizeSlice(data interface{}) interface{} {
	if data == nil {
		return nil
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice && v.IsNil() {
		return []interface{}{}
	}
	return data
}

// normalizeItems 确保分页嵌套 Items 永不为 nil — JSON 输出 [] 而非 null
func normalizeItems(data interface{}) interface{} {
	result := normalizeSlice(data)
	if result == nil {
		return []interface{}{}
	}
	return result
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(StatusSuccess, Response{
		Code:    StatusSuccess,
		Message: message,
		Data:    normalizeSlice(data),
		Success: true,
	})
}

// SuccessWithMessage 成功响应（无数据）
func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(StatusSuccess, Response{
		Code:    StatusSuccess,
		Message: message,
		Success: true,
	})
}

// BadRequest 请求参数错误（使用通用 INVALID_INPUT 码）
func BadRequest(c *gin.Context, details string) {
	FailWithCode(c, StatusBadRequest, i18n.GetMessage(c, "error.invalid_input"), "INVALID_INPUT", "error.invalid_input", details)
}

// BadRequestWithCode 请求参数错误（指定 error_code 和 message_key）
func BadRequestWithCode(c *gin.Context, code, key, details string) {
	FailWithCode(c, StatusBadRequest, i18n.GetMessage(c, key), code, key, details)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	FailWithCode(c, StatusUnauthorized, i18n.GetMessage(c, "error.unauthorized"), "UNAUTHORIZED", "error.unauthorized", message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	FailWithCode(c, StatusForbidden, i18n.GetMessage(c, "error.forbidden"), "FORBIDDEN", "error.forbidden", message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	FailWithCode(c, StatusNotFound, i18n.GetMessage(c, "error.not_found"), "NOT_FOUND", "error.not_found", message)
}

// ServerError 服务器错误
func ServerError(c *gin.Context, message string) {
	FailWithCode(c, StatusServerError, i18n.GetMessage(c, "error.internal_error"), "INTERNAL_ERROR", "error.internal_error", "")
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Success: false,
	})
}

// FailWithCode 失败响应（带错误码）
func FailWithCode(c *gin.Context, code int, message string, errorCode string, messageKey string, details string) {
	resp := Response{
		Code:       code,
		Message:    message,
		Success:    false,
		ErrorCode:  errorCode,
		MessageKey: messageKey,
	}
	if details != "" {
		resp.Details = details
	}
	c.JSON(code, resp)
}

// Error 从 error 构建统一错误响应
func Error(c *gin.Context, err error) {
	appErr := FromError(err)
	ErrorWithCode(c, appErr)
}

// FromError 将 error 映射为 AppError
func FromError(err error) *apperr.AppError {
	return apperr.FromError(err)
}

// ErrorWithCode 直接使用 AppError 构建响应
func ErrorWithCode(c *gin.Context, appErr *apperr.AppError) {
	if appErr == nil {
		FailWithCode(c, 500, i18n.GetMessage(c, "error.internal_error"), "INTERNAL_ERROR", "error.internal_error", "")
		return
	}
	FailWithCode(c, appErr.Status, i18n.GetMessage(c, appErr.MessageKey), appErr.Code, appErr.MessageKey, "")
}

// PageResult 分页结果
type PageResult struct {
	Items      interface{} `json:"items"`      // 分页数据
	Total      int64       `json:"total"`      // 总记录数
	Page       int         `json:"page"`       // 当前页码
	Size       int         `json:"size"`       // 每页大小
	TotalPages int         `json:"totalPages"` // 总页数
}

// KeysetResult keyset 游标分页结果
type KeysetResult struct {
	Items       interface{} `json:"items"`                   // 分页数据
	HasMore     bool        `json:"has_more"`                // 是否还有更多数据
	NextAfter   string      `json:"next_after,omitempty"`    // 下一页游标（时间戳）
	NextAfterID string      `json:"next_after_id,omitempty"` // 下一页游标（ID）
}

// SuccessWithPage 返回分页数据
func SuccessWithPage(c *gin.Context, items interface{}, total int64, page, size int, totalPages int) {
	Success(c, PageResult{
		Items:      normalizeItems(items),
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, "获取成功")
}

// SuccessWithKeyset 返回 keyset 游标分页数据
func SuccessWithKeyset(c *gin.Context, items interface{}, hasMore bool) {
	SuccessWithKeysetFull(c, items, hasMore, "", "")
}

// SuccessWithKeysetFull 返回 keyset 游标分页数据（含下一页游标）
func SuccessWithKeysetFull(c *gin.Context, items interface{}, hasMore bool, nextAfter, nextAfterID string) {
	Success(c, KeysetResult{
		Items:       normalizeItems(items),
		HasMore:     hasMore,
		NextAfter:   nextAfter,
		NextAfterID: nextAfterID,
	}, "获取成功")
}
