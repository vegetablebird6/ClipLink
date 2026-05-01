package repository

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// SyncEventRepository 同步事件仓库接口
type SyncEventRepository interface {
	// Save 保存同步事件
	Save(event *model.SyncEvent) error

	// FindByChannel 查找通道下的同步事件（keyset 游标分页）
	FindByChannel(channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error)

	// Count 统计通道下的同步事件数量
	Count(channelID string) (int64, error)
}
