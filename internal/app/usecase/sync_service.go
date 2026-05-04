package usecase

import (
	"context"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

type syncService struct {
	syncEventRepo repository.SyncEventRepository
	deviceRepo    repository.DeviceRepository
}

func NewSyncService(syncEventRepo repository.SyncEventRepository, deviceRepo repository.DeviceRepository) service.SyncService {
	return &syncService{
		syncEventRepo: syncEventRepo,
		deviceRepo:    deviceRepo,
	}
}

func (s *syncService) GetSyncHistory(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*service.SyncEventOutput, error) {
	events, err := s.syncEventRepo.FindByChannel(ctx, channelID, afterCreatedAt, afterID, limit)
	if err != nil {
		return nil, err
	}
	return toSyncEventOutputs(events), nil
}

func (s *syncService) LogSyncAction(ctx context.Context, actorDeviceID, channelID, content string) error {
	if actorDeviceID == "" {
		return model.ErrInvalidInput
	}

	device, err := s.deviceRepo.FindByIDAndChannel(ctx, actorDeviceID, channelID)
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

	return s.syncEventRepo.Save(ctx, event)
}

// --- model → output converters ---

func toSyncEventOutputs(events []*model.SyncEvent) []*service.SyncEventOutput {
	result := make([]*service.SyncEventOutput, 0, len(events))
	for _, e := range events {
		result = append(result, &service.SyncEventOutput{
			ID:              e.ID,
			ChannelID:       e.ChannelID,
			Action:          e.Action,
			TargetType:      e.TargetType,
			TargetID:        e.TargetID,
			Content:         e.Content,
			Summary:         e.Summary,
			ActorDeviceID:   e.ActorDeviceID,
			ActorDeviceName: e.ActorDeviceName,
			ActorDeviceType: e.ActorDeviceType,
			CreatedAt:       e.CreatedAt,
		})
	}
	return result
}
