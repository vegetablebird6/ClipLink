package repository

import (
	"context"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// SyncEventRepository 同步事件仓库接口
type SyncEventRepository interface {
	// Save 保存同步事件
	Save(ctx context.Context, event *model.SyncEvent) error

	// FindByChannel 查找通道下的同步事件（keyset 游标分页）
	FindByChannel(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error)

	// Count 统计通道下的同步事件数量
	Count(ctx context.Context, channelID string) (int64, error)
}
