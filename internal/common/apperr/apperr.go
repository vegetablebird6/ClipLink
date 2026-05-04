package apperr

import (
	"errors"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// AppError 应用级结构化错误
type AppError struct {
	Status     int
	Code       string
	MessageKey string
	Cause      error
}

func (e *AppError) Error() string { return e.MessageKey }
func (e *AppError) Unwrap() error { return e.Cause }

// New 创建 AppError
func New(status int, code, messageKey string) *AppError {
	return &AppError{Status: status, Code: code, MessageKey: messageKey}
}

// WithCause 附加原始错误原因
func (e *AppError) WithCause(cause error) *AppError {
	return &AppError{
		Status:     e.Status,
		Code:       e.Code,
		MessageKey: e.MessageKey,
		Cause:      cause,
	}
}

// 预定义错误
var (
	ErrInvalidInput       = New(400, "INVALID_INPUT", "error.invalid_input")
	ErrChannelIDRequired  = New(400, "CHANNEL_ID_REQUIRED", "error.channel_id_required")
	ErrInvalidChannelID   = New(400, "INVALID_CHANNEL_ID", "error.invalid_channel_id")
	ErrInvalidDeviceType  = New(400, "INVALID_DEVICE_TYPE", "error.invalid_device_type")
	ErrInvalidContentType = New(400, "INVALID_CLIPBOARD_TYPE", "error.invalid_clipboard_type")
	ErrInvalidContentFmt  = New(400, "INVALID_CONTENT_FORMAT", "error.invalid_content_format")
	ErrDeviceIDRequired   = New(400, "DEVICE_ID_REQUIRED", "error.device_id_required")
	ErrKeywordRequired    = New(400, "SEARCH_KEYWORD_REQUIRED", "error.search_keyword_required")
	ErrChannelNotFound    = New(404, "CHANNEL_NOT_FOUND", "error.channel_not_found")
	ErrDeviceNotFound     = New(404, "DEVICE_NOT_FOUND", "error.device_not_found")
	ErrClipboardNotFound  = New(404, "CLIPBOARD_NOT_FOUND", "error.clipboard_not_found")
	ErrUnauthorized       = New(401, "UNAUTHORIZED", "error.unauthorized")
	ErrForbidden          = New(403, "FORBIDDEN", "error.forbidden")
	ErrNotFound           = New(404, "NOT_FOUND", "error.not_found")
	ErrInternal           = New(500, "INTERNAL_ERROR", "error.internal_error")
)

// FromError 将 domain error 映射为 AppError
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, model.ErrInvalidInput):
		return ErrInvalidInput
	case errors.Is(err, model.ErrChannelNotFound):
		return ErrChannelNotFound
	case errors.Is(err, model.ErrDeviceNotFound):
		return ErrDeviceNotFound
	case errors.Is(err, model.ErrClipboardNotFound):
		return ErrClipboardNotFound
	case errors.Is(err, model.ErrUnauthorized):
		return ErrUnauthorized
	default:
		return ErrInternal
	}
}
