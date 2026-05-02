package service

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// UpdateClipboardInput 部分更新输入，只更新非 nil 字段
type UpdateClipboardInput struct {
	Title         *string
	Content       *string
	Type          *string
	DeviceType    *string
	ContentHTML   *string
	ContentFormat *string
}

// ClipboardService 剪贴板服务接口
type ClipboardService interface {
	// SaveClipboard 保存剪贴板项目
	SaveClipboard(title, content, contentType, deviceID, deviceType, channelID string, cleanDuplicates bool, contentHTML, contentFormat string) (*model.ClipboardItem, error)

	// GetLatestClipboard 获取最新的剪贴板项目
	GetLatestClipboard(channelID string, limit int) ([]*model.ClipboardItem, error)

	// GetClipboardItem 获取剪贴板项目
	GetClipboardItem(id string, channelID string) (*model.ClipboardItem, error)

	// GetClipboardHistory 获取剪贴板历史记录（keyset 游标分页）
	GetClipboardHistory(channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// DeleteClipboard 删除剪贴板项目（actorDeviceID 为执行删除操作的设备）
	DeleteClipboard(id string, channelID string, actorDeviceID string) error

	// UpdateClipboard 更新剪贴板项目（部分更新，nil 字段不更新）
	UpdateClipboard(id, channelID, actorDeviceID string, input *UpdateClipboardInput) (*model.ClipboardItem, error)

	// ToggleFavorite 切换收藏状态（actorDeviceID 为执行操作的设备）
	ToggleFavorite(id string, isFavorite bool, channelID string, actorDeviceID string) (*model.ClipboardItem, error)

	// GetFavoriteClipboard 获取收藏的剪贴板项目
	GetFavoriteClipboard(channelID string, limit int) ([]*model.ClipboardItem, error)

	// GetClipboardByType 按内容类型获取剪贴板历史记录（keyset 游标分页）
	GetClipboardByType(contentType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// GetClipboardByDeviceType 按设备类型获取剪贴板历史记录（keyset 游标分页）
	GetClipboardByDeviceType(deviceType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// GetClipboardByTypeAndDeviceType 同时按内容类型和设备类型获取剪贴板历史记录（keyset 游标分页）
	GetClipboardByTypeAndDeviceType(contentType, deviceType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// SearchClipboard 按关键词搜索剪贴板项目（offset 分页）
	SearchClipboard(keyword, channelID string, page, size int) (items []*model.ClipboardItem, total int64, totalPages int, err error)

	// CleanupDuplicateContents 清理重复内容
	CleanupDuplicateContents(channelID string) (int64, error)
}
