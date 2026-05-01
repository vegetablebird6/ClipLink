package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/common/response"
	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/model"
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

	// 绑定请求体 - 适配前端发送的字段格式
	var req struct {
		Title           string `json:"title"`
		Content         string `json:"content" binding:"required"`
		Type            string `json:"type" binding:"required"`
		DeviceID        string `json:"device_id" binding:"required"`
		DeviceType      string `json:"device_type" binding:"required"`
		CleanDuplicates bool   `json:"clean_duplicates"`
		ContentHTML     string `json:"content_html"`
		ContentFormat   string `json:"content_format"`
	}

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

	// 保存剪贴板内容
	item, err := c.clipboardService.SaveClipboard(
		req.Title,
		req.Content,
		req.Type,
		req.DeviceID,
		req.DeviceType,
		channelID,
		req.CleanDuplicates,
		req.ContentHTML,
		req.ContentFormat,
	)

	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, item, "保存成功")
}

// GetLatestClipboard 获取最新剪贴板内容
func (c *ClipboardController) GetLatestClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	// 获取查询参数
	limitStr := ctx.DefaultQuery("limit", "1") // 默认只返回1条，针对 /current 路径
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 1
	}

	// 获取最新剪贴板内容
	items, err := c.clipboardService.GetLatestClipboard(channelID, limit)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	// 如果没有找到记录，返回空数组
	if len(items) == 0 {
		response.Success(ctx, []*model.ClipboardItem{}, "获取成功")
		return
	}

	// 针对 /current 路径，只返回第一个项目而不是数组
	if ctx.Request.URL.Path == "/api/clipboard/current" || ctx.Request.URL.Path == "/clipboard/current" {
		response.Success(ctx, items[0], "获取成功")
		return
	}

	response.Success(ctx, items, "获取成功")
}

// GetClipboardItem 获取特定剪贴板项目
func (c *ClipboardController) GetClipboardItem(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	// 获取剪贴板项目
	item, err := c.clipboardService.GetClipboardItem(itemID, channelID)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	if item == nil {
		response.NotFound(ctx, "clipboard item not found")
		return
	}

	response.Success(ctx, item, "获取成功")
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
		response.ServerError(ctx, err.Error())
		return
	}

	items, hasMore := keysetHasMore(items, size)
	response.SuccessWithKeyset(ctx, items, hasMore)
}

// DeleteClipboard 删除剪贴板项目
func (c *ClipboardController) DeleteClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	// 删除剪贴板项目
	err := c.clipboardService.DeleteClipboard(itemID, channelID)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMessage(ctx, "clipboard item deleted")
}

// UpdateClipboard 更新剪贴板项目
func (c *ClipboardController) UpdateClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	// 绑定请求体 - 适配前端发送的字段格式
	var req struct {
		Title         string `json:"title"`
		Content       string `json:"content"`
		Type          string `json:"type"`
		DeviceID      string `json:"device_id"`
		DeviceType    string `json:"device_type"`
		IsFavorite    *bool  `json:"isFavorite"` // 使用指针类型，允许为空
		ContentHTML   string `json:"content_html"`
		ContentFormat string `json:"content_format"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Type != "" && !validation.IsValidClipboardType(req.Type) {
		response.BadRequest(ctx, "invalid clipboard type: "+req.Type)
		return
	}
	if req.DeviceType != "" && !validation.IsValidDeviceType(req.DeviceType) {
		response.BadRequest(ctx, "invalid device type: "+req.DeviceType)
		return
	}
	if !validation.IsValidContentFormat(req.ContentFormat) {
		response.BadRequest(ctx, "invalid content format: "+req.ContentFormat)
		return
	}

	// 更新剪贴板项目
	item, err := c.clipboardService.UpdateClipboard(
		itemID,
		req.Title,
		req.Content,
		req.Type,
		req.DeviceID,
		req.DeviceType,
		channelID,
		req.ContentHTML,
		req.ContentFormat,
	)

	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	// 如果提供了收藏状态，单独处理
	if req.IsFavorite != nil {
		item, err = c.clipboardService.ToggleFavorite(itemID, *req.IsFavorite, channelID, req.DeviceID)
		if err != nil {
			response.ServerError(ctx, err.Error())
			return
		}
	}

	response.Success(ctx, item, "更新成功")
}

// ToggleFavorite 切换收藏状态
func (c *ClipboardController) ToggleFavorite(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}
	itemID := ctx.Param("itemID")

	// 绑定请求体 - 适配前端发送的字段格式
	var req struct {
		IsFavorite bool   `json:"isFavorite" binding:"required"`
		DeviceID   string `json:"device_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 切换收藏状态
	item, err := c.clipboardService.ToggleFavorite(itemID, req.IsFavorite, channelID, req.DeviceID)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, item, "更新成功")
}

// GetFavoriteClipboard 获取收藏的剪贴板项目
func (c *ClipboardController) GetFavoriteClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	// 获取查询参数
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	// 获取收藏项目
	items, err := c.clipboardService.GetFavoriteClipboard(channelID, limit)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, items, "获取成功")
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
		response.ServerError(ctx, err.Error())
		return
	}

	items, hasMore := keysetHasMore(items, size)
	response.SuccessWithKeyset(ctx, items, hasMore)
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
		response.ServerError(ctx, err.Error())
		return
	}

	items, hasMore := keysetHasMore(items, size)
	response.SuccessWithKeyset(ctx, items, hasMore)
}

// GetCurrentClipboard 获取当前剪贴板内容（专用接口，避免路由冲突）
func (c *ClipboardController) GetCurrentClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	// 获取最新的一条剪贴板内容
	items, err := c.clipboardService.GetLatestClipboard(channelID, 1)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	// 如果没有找到记录，返回空对象
	if len(items) == 0 {
		response.Success(ctx, nil, "获取成功")
		return
	}

	// 返回第一条记录
	response.Success(ctx, items[0], "获取成功")
}

// SearchClipboard 搜索剪贴板项目
func (c *ClipboardController) SearchClipboard(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	// 获取搜索关键词
	keyword := ctx.Query("q")
	if keyword == "" {
		response.BadRequest(ctx, "搜索关键词不能为空")
		return
	}

	page, size := paginationParams(ctx, 20)

	// 执行搜索
	items, total, totalPages, err := c.clipboardService.SearchClipboard(keyword, channelID, page, size)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"items":      items,
		"total":      total,
		"page":       page,
		"size":       size,
		"totalPages": totalPages,
		"keyword":    keyword,
	}, "搜索成功")
}

// CleanupDuplicateContents 清理当前通道下已存在的重复剪贴板内容。
func (c *ClipboardController) CleanupDuplicateContents(ctx *gin.Context) {
	channelID, ok := clipboardChannelID(ctx)
	if !ok {
		return
	}

	deleted, err := c.clipboardService.CleanupDuplicateContents(channelID)
	if err != nil {
		response.ServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"deleted": deleted}, "重复内容已清理")
}
