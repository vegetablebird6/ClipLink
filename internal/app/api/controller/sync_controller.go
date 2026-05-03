package controller

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/app/api/dto"
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

	events, err := c.syncService.GetSyncHistory(channelID.(string), afterCreatedAt, afterID, limit)
	if err != nil {
		log.Printf("[sync] get history failed: %v", err)
		response.Error(ctx, err)
		return
	}

	hasMore := len(events) > limit
	if hasMore {
		events = events[:limit]
	}

	result := dto.ToSyncEventResponseList(events)
	nextAfter, nextAfterID := nextSyncCursor(result)
	response.SuccessWithKeysetFull(ctx, result, hasMore, nextAfter, nextAfterID)
}

// nextSyncCursor 从当前页最后一条的 created_at + id 计算下一页游标
func nextSyncCursor(items []*dto.SyncEventResponse) (string, string) {
	if len(items) == 0 {
		return "", ""
	}
	last := items[len(items)-1]
	return last.CreatedAt.Format(time.RFC3339Nano), strconv.FormatUint(uint64(last.ID), 10)
}

// LogSyncAction 记录同步操作
func (c *SyncController) LogSyncAction(ctx *gin.Context) {
	channelID, exists := ctx.Get("channelID")
	if !exists || channelID == nil || channelID == "" {
		response.BadRequest(ctx, "channel ID is required")
		return
	}

	var req dto.LogSyncActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "device_id and content are required")
		return
	}

	err := c.syncService.LogSyncAction(req.DeviceID, channelID.(string), req.Content)
	if err != nil {
		log.Printf("[sync] log action failed: %v", err)
		response.Error(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "sync action logged")
}
