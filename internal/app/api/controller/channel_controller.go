package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/app/api/dto"
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
	var req dto.CreateChannelRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		req.ChannelID = ""
	}
	if req.ChannelID != "" && !validation.IsValidChannelID(req.ChannelID) {
		response.BadRequest(ctx, "invalid channel ID")
		return
	}

	channel, err := c.channelService.CreateChannel(ctx.Request.Context(), req.ChannelID)
	if err != nil {
		log.Printf("[channel] create failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, dto.ToChannelResponse(channel), "频道创建成功")
}

// GetChannel 获取频道信息
func (c *ChannelController) GetChannel(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	channel, err := c.channelService.GetChannel(ctx.Request.Context(), channelID.(string))
	if err != nil {
		if err == model.ErrChannelNotFound {
			response.NotFound(ctx, "channel not found")
			return
		}
		log.Printf("[channel] get failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, dto.ToChannelResponse(channel), "获取成功")
}

// DeleteChannel 删除当前频道及其关联数据。
func (c *ChannelController) DeleteChannel(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	result, err := c.channelService.DeleteChannel(ctx.Request.Context(), channelID.(string))
	if err != nil {
		if err == model.ErrChannelNotFound {
			response.NotFound(ctx, "channel not found")
			return
		}
		log.Printf("[channel] delete failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, dto.ToChannelDeleteResponse(result), "通道已删除")
}

// VerifyChannel 验证频道存在且有效 - 适配前端POST请求格式
func (c *ChannelController) VerifyChannel(ctx *gin.Context) {
	if channelID, exists := ctx.Get("channelID"); exists && channelID != nil && channelID != "" {
		exists, err := c.channelService.VerifyChannel(ctx.Request.Context(), channelID.(string))
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

	var req dto.VerifyChannelRequest
	if err := ctx.ShouldBindJSON(&req); err == nil && req.ChannelID != "" {
		if !validation.IsValidChannelID(req.ChannelID) {
			response.BadRequest(ctx, "invalid channel ID")
			return
		}
		exists, err := c.channelService.VerifyChannel(ctx.Request.Context(), req.ChannelID)
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
