package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/common/i18n"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// ChannelAuthMiddleware 频道认证中间件
type ChannelAuthMiddleware struct {
	channelService service.ChannelService
}

// NewChannelAuthMiddleware 创建新的频道认证中间件
func NewChannelAuthMiddleware(channelService service.ChannelService) *ChannelAuthMiddleware {
	return &ChannelAuthMiddleware{
		channelService: channelService,
	}
}

// ExtractChannelFromHeader 从请求头提取频道ID并验证
func (m *ChannelAuthMiddleware) ExtractChannelFromHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		channelID := c.GetHeader("X-Channel-ID")
		if channelID == "" {
			response.BadRequestWithCode(c, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
			c.Abort()
			return
		}
		if !validation.IsValidChannelID(channelID) {
			response.BadRequestWithCode(c, "INVALID_CHANNEL_ID", "error.invalid_channel_id", "")
			c.Abort()
			return
		}

		exists, err := m.channelService.VerifyChannel(c.Request.Context(), channelID)
		if err != nil {
			log.Printf("[channel auth] verify failed: %v", err)
			response.FailWithCode(c, 500, i18n.GetMessage(c, "error.internal_error"), "INTERNAL_ERROR", "error.internal_error", "")
			c.Abort()
			return
		}

		if !exists {
			response.FailWithCode(c, 404, i18n.GetMessage(c, "error.channel_not_found"), "CHANNEL_NOT_FOUND", "error.channel_not_found", "")
			c.Abort()
			return
		}

		c.Set("channelID", channelID)
		c.Next()
	}
}

// VerifyChannel 验证频道是否存在且有效 (路径参数版)
func (m *ChannelAuthMiddleware) VerifyChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		channelID := c.Param("channelID")
		if channelID == "" {
			response.BadRequestWithCode(c, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
			c.Abort()
			return
		}
		if !validation.IsValidChannelID(channelID) {
			response.BadRequestWithCode(c, "INVALID_CHANNEL_ID", "error.invalid_channel_id", "")
			c.Abort()
			return
		}

		exists, err := m.channelService.VerifyChannel(c.Request.Context(), channelID)
		if err != nil {
			log.Printf("[channel auth] verify failed: %v", err)
			response.FailWithCode(c, 500, i18n.GetMessage(c, "error.internal_error"), "INTERNAL_ERROR", "error.internal_error", "")
			c.Abort()
			return
		}

		if !exists {
			response.FailWithCode(c, 404, i18n.GetMessage(c, "error.channel_not_found"), "CHANNEL_NOT_FOUND", "error.channel_not_found", "")
			c.Abort()
			return
		}

		c.Next()
	}
}
