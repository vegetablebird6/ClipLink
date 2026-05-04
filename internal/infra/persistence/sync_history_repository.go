package persistence

import (
	"context"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"gorm.io/gorm"
)

// syncEventRepository 同步事件仓库实现
type syncEventRepository struct {
	gdb *gorm.DB
}

// NewSyncEventRepository 创建新的同步事件仓库
func NewSyncEventRepository(gdb *gorm.DB) repository.SyncEventRepository {
	return &syncEventRepository{gdb: gdb}
}

// Save 保存同步事件
func (r *syncEventRepository) Save(ctx context.Context, event *model.SyncEvent) error {
	return r.gdb.WithContext(ctx).Create(event).Error
}

// FindByChannel 查找通道下的同步事件（keyset 游标分页）
// 查询 limit+1 条，调用方用 len(events) > limit 判断 has_more。
func (r *syncEventRepository) FindByChannel(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error) {
	var events []*model.SyncEvent
	query := r.gdb.WithContext(ctx).Where("channel_id = ?", channelID)
	if afterCreatedAt != nil && afterID != nil {
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			*afterCreatedAt, *afterCreatedAt, *afterID)
	}
	err := query.Order("created_at DESC, id DESC").Limit(limit + 1).Find(&events).Error
	return events, err
}

// Count 统计通道下的同步事件数量
func (r *syncEventRepository) Count(ctx context.Context, channelID string) (int64, error) {
	var count int64
	query := r.gdb.WithContext(ctx).Model(&model.SyncEvent{})

	if channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}

	err := query.Count(&count).Error
	return count, err
}
