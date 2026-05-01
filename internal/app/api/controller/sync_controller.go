package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// SyncController 同步控制器
type SyncController struct {
	syncService service.SyncService
}

// NewSyncController 创建新的同步控制器
func NewSyncController(syncService service.SyncService) *SyncController {
	return &SyncController{
		syncService: syncService,
	}
}

// GetSyncHistory 获取同步事件记录（keyset 游标分页）
func (c *SyncController) GetSyncHistory(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	// 获取分页参数
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// 解析 keyset 游标
	var afterCreatedAt *time.Time
	var afterID *uint
	afterStr := ctx.Query("after")
	afterIDStr := ctx.Query("after_id")
	if afterStr != "" && afterIDStr != "" {
		if t, parseErr := time.Parse(time.RFC3339Nano, afterStr); parseErr == nil {
			afterCreatedAt = &t
			if idVal, idErr := strconv.ParseUint(afterIDStr, 10, 64); idErr == nil {
				uintID := uint(idVal)
				afterID = &uintID
			}
		}
	}

	// 获取同步事件记录
	events, err := c.syncService.GetSyncHistory(channelID.(string), afterCreatedAt, afterID, limit)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	// 判断 has_more：多查了 1 条，如果取到 limit+1 条说明还有更多
	hasMore := len(events) > limit
	if hasMore {
		events = events[:limit]
	}
	response.SuccessWithKeyset(ctx, events, hasMore)
}

// LogSyncAction 记录同步操作
func (c *SyncController) LogSyncAction(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	// 绑定请求体
	var req struct {
		DeviceID string `json:"deviceId" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 记录同步操作
	err := c.syncService.LogSyncAction(req.DeviceID, channelID.(string), req.Content)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "sync action logged")
}
