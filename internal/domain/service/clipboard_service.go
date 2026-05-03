package service

import "time"

// ClipboardService 剪贴板服务接口
type ClipboardService interface {
	// CreateClipboard 创建剪贴板条目
	CreateClipboard(in CreateClipboardInput) (*ClipboardItemOutput, error)

	// GetLatestClipboard 获取最新的剪贴板条目
	GetLatestClipboard(channelID string, limit int) ([]*ClipboardItemOutput, error)

	// GetClipboardItem 获取剪贴板条目
	GetClipboardItem(id, channelID string) (*ClipboardItemOutput, error)

	// GetClipboardHistory 获取剪贴板历史记录（keyset 游标分页）
	GetClipboardHistory(channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)

	// DeleteClipboard 删除剪贴板条目
	DeleteClipboard(in DeleteClipboardInput) error

	// UpdateClipboard 更新剪贴板条目（部分更新，nil 字段不更新）
	UpdateClipboard(in UpdateClipboardInput) (*ClipboardItemOutput, error)

	// SetFavorite 设置收藏状态
	SetFavorite(in SetFavoriteInput) (*ClipboardItemOutput, error)

	// GetFavoriteClipboard 获取收藏的剪贴板条目（keyset 游标分页）
	GetFavoriteClipboard(channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)

	// GetClipboardByType 按内容类型获取剪贴板历史记录（keyset 游标分页）
	GetClipboardByType(contentType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)

	// GetClipboardByDeviceType 按设备类型获取剪贴板历史记录（keyset 游标分页）
	GetClipboardByDeviceType(deviceType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)

	// GetClipboardByTypeAndDeviceType 同时按内容类型和设备类型获取剪贴板历史记录（keyset 游标分页）
	GetClipboardByTypeAndDeviceType(contentType, deviceType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)

	// SearchClipboard 按关键词搜索剪贴板条目（offset 分页）
	SearchClipboard(keyword, channelID string, page, size int) (items []*ClipboardItemOutput, total int64, totalPages int, err error)

	// CleanupDuplicateContents 清理重复内容
	CleanupDuplicateContents(channelID string) (int64, error)
}
