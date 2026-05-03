package controller

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xiaojiu/cliplink/internal/app/api/dto"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// ClipboardController 剪贴板控制器
type ClipboardController struct {
	clipboardService service.ClipboardService
}

// NewClipboardController 创建新的剪贴板控制器
func NewClipboardController(clipboardService service.ClipboardService) *ClipboardController {
	return &ClipboardController{
		clipboardService: clipboardService,
	}
}

func clipboardChannelID(ctx *gin.Context) (string, bool) {
	channelID, exists := ctx.Get("channelID")
	if !exists {
		response.BadRequest(ctx, "channel ID is required")
		return "", false
	}
	value, ok := channelID.(string)
	if !ok || value == "" {
		response.BadRequest(ctx, "channel ID is required")
		return "", false
	}
	return value, true
}

func respondClipboardError(ctx *gin.Context, err error) {
	log.Printf("[clipboard error] %v", err)
	response.Error(ctx, err)
}

func paginationParams(ctx *gin.Context, defaultSize int) (int, int) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(ctx.DefaultQuery("size", strconv.Itoa(defaultSize)))
	if err != nil || size < 1 {
		size = defaultSize
	}
	if size > 100 {
		size = 100
	}

	return page, size
}

// keysetCursor 解析 keyset 游标分页参数：after（ISO 时间戳）和 after_id（上页最后一条 ID）
func keysetCursor(ctx *gin.Context) (*time.Time, *string) {
	afterStr := ctx.Query("after")
	afterID := ctx.Query("after_id")
	if afterStr == "" || afterID == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339Nano, afterStr)
	if err != nil {
		// 兼容 ISO 8601 无时区格式
		t, err = time.Parse("2006-01-02T15:04:05.999999999", afterStr)
		if err != nil {
			return nil, nil
		}
	}
	return &t, &afterID
}

// keysetSize 解析 keyset 分页的 size 参数
func keysetSize(ctx *gin.Context, defaultSize int) int {
	size, err := strconv.Atoi(ctx.DefaultQuery("size", strconv.Itoa(defaultSize)))
	if err != nil || size < 1 {
		size = defaultSize
	}
	if size > 100 {
		size = 100
	}
	return size
}

// keysetHasMore 判断 keyset 分页是否还有更多数据，并裁掉 size+1 条中的额外记录。
// 调用约定：repository 查 size+1 条，items 可能有 size 或 size+1 条。
func keysetHasMore[T any](items []T, size int) ([]T, bool) {
	if len(items) > size {
		return items[:size], true
	}
	return items, false
}

// SaveClipboard 保存剪贴板内容
func (c *ClipboardController) SaveClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	var req dto.CreateClipboardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if !validation.IsValidClipboardType(req.Type) {
		response.BadRequest(ctx, "invalid clipboard type: "+req.Type)
		return
	}
	if !validation.IsValidDeviceType(req.DeviceType) {
		response.BadRequest(ctx, "invalid device type: "+req.DeviceType)
		return
	}
	if !validation.IsValidContentFormat(req.ContentFormat) {
		response.BadRequest(ctx, "invalid content format: "+req.ContentFormat)
		return
	}

	item, err := c.clipboardService.CreateClipboard(service.CreateClipboardInput{
		ChannelID:       channelID,
		ActorDeviceID:   req.DeviceID,
		ActorDeviceType: req.DeviceType,
		Title:           req.Title,
		Content:         req.Content,
		Type:            req.Type,
		CleanDuplicates: req.CleanDuplicates,
		ContentHTML:     req.ContentHTML,
		ContentFormat:   req.ContentFormat,
	})

	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	response.Success(ctx, dto.ToClipboardItemResponse(item), "保存成功")
}

// GetClipboardItem 获取特定剪贴板项目
func (c *ClipboardController) GetClipboardItem(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	item, err := c.clipboardService.GetClipboardItem(itemID, channelID)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	if item == nil {
		response.NotFound(ctx, "clipboard item not found")
		return
	}

	response.Success(ctx, dto.ToClipboardItemResponse(item), "获取成功")
}

// GetClipboardHistory 获取剪贴板历史记录（keyset 游标分页）
func (c *ClipboardController) GetClipboardHistory(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	size := keysetSize(ctx, 20)
	afterCreatedAt, afterID := keysetCursor(ctx)

	items, err := c.clipboardService.GetClipboardHistory(channelID, afterCreatedAt, afterID, size)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	items, hasMore := keysetHasMore(items, size)
	nextAfter, nextAfterID := nextKeysetCursor(items)
	response.SuccessWithKeysetFull(ctx, dto.ToClipboardItemResponseList(items), hasMore, nextAfter, nextAfterID)
}

// nextKeysetCursor 从当前页最后一条的 created_at + id 计算下一页游标
func nextKeysetCursor(items []*service.ClipboardItemOutput) (string, string) {
	if len(items) == 0 {
		return "", ""
	}
	last := items[len(items)-1]
	return last.CreatedAt.Format(time.RFC3339Nano), last.ID
}

// nextFavoritesCursor 从当前页最后一条的 updated_at + id 计算下一页游标（收藏按 updated_at 排序）
func nextFavoritesCursor(items []*service.ClipboardItemOutput) (string, string) {
	if len(items) == 0 {
		return "", ""
	}
	last := items[len(items)-1]
	return last.UpdatedAt.Format(time.RFC3339Nano), last.ID
}

// DeleteClipboard 删除剪贴板项目
func (c *ClipboardController) DeleteClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	var req dto.DeleteClipboardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "device_id is required")
		return
	}

	err := c.clipboardService.DeleteClipboard(service.DeleteClipboardInput{
		ID:            itemID,
		ChannelID:     channelID,
		ActorDeviceID: req.DeviceID,
	})
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	response.SuccessWithMessage(ctx, "clipboard item deleted")
}

// UpdateClipboard 更新剪贴板项目（部分更新）
func (c *ClipboardController) UpdateClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	var req dto.UpdateClipboardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Type != nil && !validation.IsValidClipboardType(*req.Type) {
		response.BadRequest(ctx, "invalid clipboard type: "+*req.Type)
		return
	}
	if req.DeviceType != nil && !validation.IsValidDeviceType(*req.DeviceType) {
		response.BadRequest(ctx, "invalid device type: "+*req.DeviceType)
		return
	}
	if req.ContentFormat != nil && !validation.IsValidContentFormat(*req.ContentFormat) {
		response.BadRequest(ctx, "invalid content format: "+*req.ContentFormat)
		return
	}

	item, err := c.clipboardService.UpdateClipboard(service.UpdateClipboardInput{
		ID:             itemID,
		ChannelID:      channelID,
		ActorDeviceID:  req.DeviceID,
		Title:          req.Title,
		Content:        req.Content,
		Type:           req.Type,
		DeviceType:     req.DeviceType,
		ContentHTML:    req.ContentHTML,
		ContentFormat:  req.ContentFormat,
	})

	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	response.Success(ctx, dto.ToClipboardItemResponse(item), "更新成功")
}

// ToggleFavorite 切换收藏状态
func (c *ClipboardController) ToggleFavorite(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	var req dto.SetFavoriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "favorite and device_id are required")
		return
	}

	item, err := c.clipboardService.SetFavorite(service.SetFavoriteInput{
		ID:            itemID,
		ChannelID:     channelID,
		ActorDeviceID: req.DeviceID,
		Favorite:      req.Favorite,
	})
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	response.Success(ctx, dto.ToClipboardItemResponse(item), "更新成功")
}

// GetFavoriteClipboard 获取收藏的剪贴板项目（keyset 游标分页）
func (c *ClipboardController) GetFavoriteClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	size := keysetSize(ctx, 20)
	afterCreatedAt, afterID := keysetCursor(ctx)

	items, err := c.clipboardService.GetFavoriteClipboard(channelID, afterCreatedAt, afterID, size)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	items, hasMore := keysetHasMore(items, size)
	nextAfter, nextAfterID := nextFavoritesCursor(items)
	response.SuccessWithKeysetFull(ctx, dto.ToClipboardItemResponseList(items), hasMore, nextAfter, nextAfterID)
}

// GetClipboardByType 按类型获取剪贴板项目（keyset 游标分页）
func (c *ClipboardController) GetClipboardByType(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	clipType := ctx.Param("type")
	if !validation.IsValidClipboardType(clipType) {
		response.BadRequest(ctx, "invalid clipboard type: "+clipType)
		return
	}

	size := keysetSize(ctx, 20)
	afterCreatedAt, afterID := keysetCursor(ctx)

	items, err := c.clipboardService.GetClipboardByType(clipType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	items, hasMore := keysetHasMore(items, size)
	nextAfter, nextAfterID := nextKeysetCursor(items)
	response.SuccessWithKeysetFull(ctx, dto.ToClipboardItemResponseList(items), hasMore, nextAfter, nextAfterID)
}

// GetClipboardByDeviceType 按设备类型获取剪贴板项目（keyset 游标分页）
func (c *ClipboardController) GetClipboardByDeviceType(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	deviceType := ctx.Param("deviceType")
	if !validation.IsValidDeviceType(deviceType) {
		response.BadRequest(ctx, "invalid device type: "+deviceType)
		return
	}

	size := keysetSize(ctx, 20)
	afterCreatedAt, afterID := keysetCursor(ctx)

	items, err := c.clipboardService.GetClipboardByDeviceType(deviceType, channelID, afterCreatedAt, afterID, size)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	items, hasMore := keysetHasMore(items, size)
	nextAfter, nextAfterID := nextKeysetCursor(items)
	response.SuccessWithKeysetFull(ctx, dto.ToClipboardItemResponseList(items), hasMore, nextAfter, nextAfterID)
}

// GetCurrentClipboard 获取当前剪贴板内容
func (c *ClipboardController) GetCurrentClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	items, err := c.clipboardService.GetLatestClipboard(channelID, 1)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	if len(items) == 0 {
		response.Success(ctx, nil, "获取成功")
		return
	}

	response.Success(ctx, dto.ToClipboardItemResponse(items[0]), "获取成功")
}

// SearchClipboard 搜索剪贴板项目
func (c *ClipboardController) SearchClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	keyword := ctx.Query("q")
	if keyword == "" {
		response.BadRequest(ctx, "搜索关键词不能为空")
		return
	}

	page, size := paginationParams(ctx, 20)

	items, total, totalPages, err := c.clipboardService.SearchClipboard(keyword, channelID, page, size)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	response.SuccessWithPage(ctx, dto.ToClipboardItemResponseList(items), total, page, size, totalPages)
}

// CleanupDuplicateContents 清理当前通道下已存在的重复剪贴板内容。
func (c *ClipboardController) CleanupDuplicateContents(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	deleted, err := c.clipboardService.CleanupDuplicateContents(channelID)
	if err != nil {
		respondClipboardError(ctx, err)
		return
	}

	response.Success(ctx, struct {
		Deleted int64 `json:"deleted"`
	}{Deleted: deleted}, "重复内容已清理")
}
