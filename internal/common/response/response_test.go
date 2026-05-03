package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xiaojiu/cliplink/internal/domain/model"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestError_WithClipboardNotFound_ZhCN(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	Error(c, model.ErrClipboardNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"success":false`)
	assert.Contains(t, w.Body.String(), `"error_code":"CLIPBOARD_NOT_FOUND"`)
	assert.Contains(t, w.Body.String(), `"message_key":"error.clipboard_not_found"`)
	assert.Contains(t, w.Body.String(), "剪贴板项不存在")
	assert.NotContains(t, w.Body.String(), "clipboard item not found")
}

func TestError_WithClipboardNotFound_EnUS(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "en-US")

	Error(c, model.ErrClipboardNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":"CLIPBOARD_NOT_FOUND"`)
	assert.Contains(t, w.Body.String(), "Clipboard item not found")
	assert.NotContains(t, w.Body.String(), "剪贴板项不存在")
}

func TestError_WithDeviceNotFound_ZhCN(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	Error(c, model.ErrDeviceNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":"DEVICE_NOT_FOUND"`)
	assert.Contains(t, w.Body.String(), "设备不存在")
}

func TestError_WithInvalidInput_ZhCN(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	Error(c, model.ErrInvalidInput)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":"INVALID_INPUT"`)
	assert.Contains(t, w.Body.String(), "参数错误")
}

func TestError_WithChannelNotFound_ZhCN(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	Error(c, model.ErrChannelNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":"CHANNEL_NOT_FOUND"`)
	assert.Contains(t, w.Body.String(), "通道不存在")
}

func TestError_WithUnknownError_NotExposed(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	unknownErr := errors.New("some internal database error message")
	Error(c, unknownErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":"INTERNAL_ERROR"`)
	assert.NotContains(t, w.Body.String(), "some internal database error message")
	assert.NotContains(t, w.Body.String(), "database")
}

func TestError_WithNilError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	Error(c, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":"INTERNAL_ERROR"`)
}

func TestError_DefaultLangIsZhCN(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	// no Accept-Language header

	Error(c, model.ErrDeviceNotFound)

	assert.Contains(t, w.Body.String(), "设备不存在")
}
