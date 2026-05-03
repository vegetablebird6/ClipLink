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
	ErrInvalidInput      = New(400, "INVALID_INPUT", "error.invalid_input")
	ErrChannelNotFound   = New(404, "CHANNEL_NOT_FOUND", "error.channel_not_found")
	ErrDeviceNotFound    = New(404, "DEVICE_NOT_FOUND", "error.device_not_found")
	ErrClipboardNotFound = New(404, "CLIPBOARD_NOT_FOUND", "error.clipboard_not_found")
	ErrUnauthorized      = New(401, "UNAUTHORIZED", "error.unauthorized")
	ErrForbidden         = New(403, "FORBIDDEN", "error.forbidden")
	ErrInternal          = New(500, "INTERNAL_ERROR", "error.internal_error")
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
