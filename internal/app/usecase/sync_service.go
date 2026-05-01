package usecase

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// syncService 同步服务实现
type syncService struct {
	syncEventRepo repository.SyncEventRepository
}

// NewSyncService 创建新的同步服务
func NewSyncService(syncEventRepo repository.SyncEventRepository) service.SyncService {
	return &syncService{
		syncEventRepo: syncEventRepo,
	}
}

// GetSyncHistory 获取同步事件记录（keyset 游标分页）
func (s *syncService) GetSyncHistory(channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error) {
	return s.syncEventRepo.FindByChannel(channelID, afterCreatedAt, afterID, limit)
}

// LogSyncAction 记录同步操作
func (s *syncService) LogSyncAction(deviceID, channelID, content string) error {
	event := &model.SyncEvent{
		Action:     model.ActionSync,
		Content:    content,
		DeviceID:   deviceID,
		ChannelID:  channelID,
		TargetType: model.TargetTypeClipboard,
		CreatedAt:  time.Now(),
	}

	return s.syncEventRepo.Save(event)
}
