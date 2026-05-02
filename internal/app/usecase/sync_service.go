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
	deviceRepo    repository.DeviceRepository
}

// NewSyncService 创建新的同步服务
func NewSyncService(syncEventRepo repository.SyncEventRepository, deviceRepo repository.DeviceRepository) service.SyncService {
	return &syncService{
		syncEventRepo: syncEventRepo,
		deviceRepo:    deviceRepo,
	}
}

// GetSyncHistory 获取同步事件记录（keyset 游标分页）
func (s *syncService) GetSyncHistory(channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*model.SyncEvent, error) {
	return s.syncEventRepo.FindByChannel(channelID, afterCreatedAt, afterID, limit)
}

// LogSyncAction 记录同步操作（actorDeviceID 必须是已注册且属于 channel 的设备）
func (s *syncService) LogSyncAction(actorDeviceID, channelID, content string) error {
	if actorDeviceID == "" {
		return model.ErrInvalidInput
	}

	device, err := s.deviceRepo.FindByIDAndChannel(actorDeviceID, channelID)
	if err != nil || device == nil {
		return model.ErrInvalidInput
	}

	event := &model.SyncEvent{
		Action:          model.ActionSync,
		Content:         content,
		ChannelID:       channelID,
		TargetType:      model.TargetTypeClipboard,
		ActorDeviceID:   device.ID,
		ActorDeviceName: device.Name,
		ActorDeviceType: device.Type,
		CreatedAt:       time.Now(),
	}

	return s.syncEventRepo.Save(event)
}
