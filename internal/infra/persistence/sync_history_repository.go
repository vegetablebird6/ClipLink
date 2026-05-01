package persistence

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/infra/db"
)

// syncEventRepository 同步事件仓库实现
type syncEventRepository struct{}

// NewSyncEventRepository 创建新的同步事件仓库
func NewSyncEventRepository() repository.SyncEventRepository {
	return &syncEventRepository{}
}

// Save 保存同步事件
func (r *syncEventRepository) Save(event *model.SyncEvent) error {
	return db.GetDB().Create(event).Error
}

// FindByChannel 查找通道下的同步事件（keyset 游标分页）
// 查询 limit+1 条，调用方用 len(events) > limit 判断 has_more。
func (r *syncEventRepository) FindByChannel(channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error) {
	var events []*model.SyncEvent
	query := db.GetDB().Where("channel_id = ?", channelID)
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("created_at DESC, id DESC").Limit(limit + 1).Find(&events).Error
	return events, err
}

// Count 统计通道下的同步事件数量
func (r *syncEventRepository) Count(channelID string) (int64, error) {
	var count int64
	query := db.GetDB().Model(&model.SyncEvent{})

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Count(&count).Error
	return count, err
}
