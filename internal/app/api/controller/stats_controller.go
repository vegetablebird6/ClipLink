package controller

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/xiaojiu/cliplink/internal/app/api/dto"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// StatsController 统计控制器
type StatsController struct {
	statsService   service.StatsService
	channelService service.ChannelService
}

// NewStatsController 创建新的统计控制器
func NewStatsController(statsService service.StatsService, channelService service.ChannelService) *StatsController {
	return &StatsController{
		statsService:   statsService,
		channelService: channelService,
	}
}

// GetChannelStats 获取通道统计数据
func (c *StatsController) GetChannelStats(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequestWithCode(ctx, "CHANNEL_ID_REQUIRED", "error.channel_id_required", "")
		return
	}

	channel, err := c.channelService.GetChannel(ctx.Request.Context(), channelID.(string))
	if err != nil {
		log.Printf("[stats] get channel failed: %v", err)
		response.Error(ctx, err)
		return
	}

	stats, err := c.statsService.GetChannelStats(ctx.Request.Context(), channelID.(string))
	if err != nil {
		log.Printf("[stats] get stats failed: %v", err)
		response.Error(ctx, err)
		return
	}

	if stats == nil {
		response.NotFound(ctx, "channel not found")
		return
	}

	response.Success(ctx, dto.NewChannelStatsResponse(stats, channel.CreatedAt), "获取成功")
}
