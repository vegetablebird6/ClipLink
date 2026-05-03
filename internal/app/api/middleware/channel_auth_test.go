package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/xiaojiu/cliplink/internal/common/response"
)

func TestExtractChannelFromHeader_MissingHeader_ZhCN(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", (&ChannelAuthMiddleware{}).ExtractChannelFromHeader(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// no X-Channel-ID header, no Accept-Language (default zh-CN)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "INVALID_INPUT", resp.ErrorCode)
	assert.Equal(t, "error.channel_id_required", resp.MessageKey)
	assert.Equal(t, "通道 ID 不能为空", resp.Message)
}

func TestExtractChannelFromHeader_MissingHeader_EnUS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", (&ChannelAuthMiddleware{}).ExtractChannelFromHeader(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Language", "en-US")
	// no X-Channel-ID header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_INPUT", resp.ErrorCode)
	assert.Equal(t, "Channel ID is required", resp.Message)
	assert.NotContains(t, w.Body.String(), "通道 ID")
}

func TestExtractChannelFromHeader_InvalidID_ZhCN(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", (&ChannelAuthMiddleware{}).ExtractChannelFromHeader(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Channel-ID", "invalid!@#")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_INPUT", resp.ErrorCode)
	assert.Equal(t, "无效的通道 ID", resp.Message)
}
