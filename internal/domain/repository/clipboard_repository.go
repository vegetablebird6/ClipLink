package repository

import (
	"context"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// ClipboardRepository 剪贴板仓库接口
type ClipboardRepository interface {
	// Save 保存剪贴板项目
	Save(ctx context.Context, item *model.ClipboardItem) error

	// FindByID 通过ID查找剪贴板项目
	FindByID(ctx context.Context, id, channelID string) (*model.ClipboardItem, error)

	// FindLatest 获取最新的剪贴板项目
	FindLatest(ctx context.Context, channelID string, limit int) ([]*model.ClipboardItem, error)

	// FindWithKeyset 分页获取剪贴板项目（keyset 游标分页）
	FindWithKeyset(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// FindByType 按类型查找剪贴板项目（keyset 游标分页）
	FindByType(ctx context.Context, contentType, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// FindByDeviceType 按设备类型查找剪贴板项目（keyset 游标分页）
	FindByDeviceType(ctx context.Context, deviceType, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// FindByTypeAndDeviceType 同时按内容类型和设备类型查找剪贴板项目（keyset 游标分页）
	FindByTypeAndDeviceType(ctx context.Context, contentType, deviceType, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// FindFavorites 查找收藏的剪贴板项目（keyset 游标分页）
	FindFavorites(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*model.ClipboardItem, error)

	// Update 更新剪贴板项目
	Update(ctx context.Context, id, channelID string, updates map[string]interface{}) error

	// Delete 删除剪贴板项目
	Delete(ctx context.Context, id, channelID string) error

	// DeleteByContentHash 基于内容哈希删除同通道下的重复项，保留指定项目。
	DeleteByContentHash(ctx context.Context, channelID, contentHash, keepID string) (int64, error)

	// CleanupDuplicateContents 清理同通道下已存在的重复内容，保留每组最新项目。
	CleanupDuplicateContents(ctx context.Context, channelID string) (int64, error)

	// Count 统计剪贴板项目数量
	Count(ctx context.Context, channelID string) (int64, error)

	// CountByType 按类型统计剪贴板项目数量
	CountByType(ctx context.Context, contentType, channelID string) (int64, error)

	// SearchByKeyword 按关键词搜索剪贴板项目（offset 分页，结果按相关度排序）
	SearchByKeyword(ctx context.Context, keyword, channelID string, page, size int) ([]*model.ClipboardItem, int64, int, error)
}
