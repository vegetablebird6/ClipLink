package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

const orphanDeviceRetention = 30 * 24 * time.Hour

type channelService struct {
	channelRepo   repository.ChannelRepository
	clipboardRepo repository.ClipboardRepository
	deviceRepo    repository.DeviceRepository
}

func NewChannelService(
	channelRepo repository.ChannelRepository,
	clipboardRepo repository.ClipboardRepository,
	deviceRepo repository.DeviceRepository,
) service.ChannelService {
	return &channelService{
		channelRepo:   channelRepo,
		clipboardRepo: clipboardRepo,
		deviceRepo:    deviceRepo,
	}
}

func (s *channelService) CreateChannel(ctx context.Context, channelID string) (*service.ChannelOutput, error) {
	id := channelID
	if id == "" {
		id = uuid.New().String()
	}

	if channelID != "" {
		exists, err := s.channelRepo.Exists(ctx, channelID)
		if err != nil {
			return nil, err
		}
		if exists {
			ch, err := s.channelRepo.FindByID(ctx, channelID)
			if err != nil {
				return nil, err
			}
			return toChannelOutput(ch), nil
		}
	}

	channel := &model.Channel{
		ID:        id,
		CreatedAt: time.Now(),
	}

	if err := s.channelRepo.Save(ctx, channel); err != nil {
		return nil, err
	}

	return toChannelOutput(channel), nil
}

func (s *channelService) GetChannel(ctx context.Context, channelID string) (*service.ChannelOutput, error) {
	ch, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	return toChannelOutput(ch), nil
}

func (s *channelService) ChannelExists(ctx context.Context, channelID string) (bool, error) {
	return s.channelRepo.Exists(ctx, channelID)
}

func (s *channelService) VerifyChannel(ctx context.Context, channelID string) (bool, error) {
	exists, err := s.channelRepo.Exists(ctx, channelID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *channelService) DeleteChannel(ctx context.Context, channelID string) (*service.ChannelDeleteOutput, error) {
	exists, err := s.channelRepo.Exists(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, model.ErrChannelNotFound
	}

	result, err := s.channelRepo.Delete(ctx, channelID, time.Now().Add(-orphanDeviceRetention))
	if err != nil {
		return nil, err
	}

	return &service.ChannelDeleteOutput{
		ChannelID:             result.ChannelID,
		ClipboardItemsDeleted: result.ClipboardItemsDeleted,
		SyncEventsDeleted:     result.SyncEventsDeleted,
		DeviceLinksDeleted:    result.DeviceLinksDeleted,
		OrphanDevicesDeleted:  result.OrphanDevicesDeleted,
	}, nil
}

func toChannelOutput(channel *model.Channel) *service.ChannelOutput {
	if channel == nil {
		return nil
	}
	return &service.ChannelOutput{
		ID:        channel.ID,
		CreatedAt: channel.CreatedAt,
	}
}
