package service

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// SyncService 同步服务接口
type SyncService interface {
	// GetSyncHistory 获取同步事件记录（keyset 游标分页）
	GetSyncHistory(channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error)

	// LogSyncAction 记录同步操作
	LogSyncAction(deviceID, channelID, content string) error
}
