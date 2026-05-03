package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// ChannelController 频道控制器
type ChannelController struct {
	channelService service.ChannelService
}

// NewChannelController 创建新的频道控制器
func NewChannelController(channelService service.ChannelService) *ChannelController {
	return &ChannelController{
		channelService: channelService,
	}
}

// CreateChannel 创建新频道
func (c *ChannelController) CreateChannel(ctx *gin.Context) {
	// 绑定请求体 - 适配前端API格式
	var req struct {
		ChannelID string `json:"channel_id"` // 允许客户端指定channelID
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果没有请求体或请求体解析错误，使用空值创建随机频道ID
		req.ChannelID = ""
	}
	if req.ChannelID != "" && !validation.IsValidChannelID(req.ChannelID) {
		response.BadRequest(ctx, "invalid channel ID")
		return
	}

	// 创建频道（使用指定的ID或生成随机ID）
	channel, err := c.channelService.CreateChannel(req.ChannelID)
	if err != nil {
		log.Printf("[channel] create failed: %v", err)
		response.Error(ctx, err)
		return
	}

	// 返回新创建的频道信息
	response.Success(ctx, gin.H{
		"id":         channel.ID,
		"created_at": channel.CreatedAt,
	}, "频道创建成功")
}

// GetChannel 获取频道信息
func (c *ChannelController) GetChannel(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	// 获取频道
	channel, err := c.channelService.GetChannel(channelID.(string))
	if err != nil {
		if err == model.ErrChannelNotFound {
			response.NotFound(ctx, "channel not found")
			return
		}
		log.Printf("[channel] get failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, channel, "获取成功")
}

// DeleteChannel 删除当前频道及其关联数据。
func (c *ChannelController) DeleteChannel(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	result, err := c.channelService.DeleteChannel(channelID.(string))
	if err != nil {
		if err == model.ErrChannelNotFound {
			response.NotFound(ctx, "channel not found")
			return
		}
		log.Printf("[channel] delete failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, result, "通道已删除")
}

// GetChannelStats 获取频道统计信息
func (c *ChannelController) GetChannelStats(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	// 获取统计信息
	stats, err := c.channelService.GetChannelStats(channelID.(string))
	if err != nil {
		if err == model.ErrChannelNotFound {
			response.NotFound(ctx, "channel not found")
			return
		}
		log.Printf("[channel] get stats failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, stats, "获取成功")
}

// VerifyChannel 验证频道存在且有效 - 适配前端POST请求格式
func (c *ChannelController) VerifyChannel(ctx *gin.Context) {
	// 优先从header注入的ctx.Get("channelID")获取
	if channelID, exists := ctx.Get("channelID"); exists && channelID != nil && channelID != "" {
		exists, err := c.channelService.VerifyChannel(channelID.(string))
		if err != nil {
			log.Printf("[channel] verify failed: %v", err)
			response.Error(ctx, err)
			return
		}
		if !exists {
			response.NotFound(ctx, "channel not found")
			return
		}
		response.Success(ctx, nil, "频道有效")
		return
	}

	// 其次尝试POST body
	var req struct {
		ChannelID string `json:"channel_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err == nil && req.ChannelID != "" {
		if !validation.IsValidChannelID(req.ChannelID) {
			response.BadRequest(ctx, "invalid channel ID")
			return
		}
		exists, err := c.channelService.VerifyChannel(req.ChannelID)
		if err != nil {
			log.Printf("[channel] verify failed: %v", err)
			response.Error(ctx, err)
			return
		}
		if !exists {
			response.NotFound(ctx, "channel not found")
			return
		}
		response.Success(ctx, nil, "频道有效")
		return
	}

	response.BadRequest(ctx, "channel ID is required")
}
