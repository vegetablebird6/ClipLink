package usecase

import (
	"context"

	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

type statsService struct {
	deviceRepo    repository.DeviceRepository
	clipboardRepo repository.ClipboardRepository
	channelRepo   repository.ChannelRepository
	syncEventRepo repository.SyncEventRepository
}

func NewStatsService(
	deviceRepo repository.DeviceRepository,
	clipboardRepo repository.ClipboardRepository,
	channelRepo repository.ChannelRepository,
	syncEventRepo repository.SyncEventRepository,
) service.StatsService {
	return &statsService{
		deviceRepo:    deviceRepo,
		clipboardRepo: clipboardRepo,
		channelRepo:   channelRepo,
		syncEventRepo: syncEventRepo,
	}
}

func (s *statsService) GetChannelStats(ctx context.Context, channelID string) (*service.StatsOutput, error) {
	exists, err := s.channelRepo.Exists(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	clipboardCount, err := s.clipboardRepo.Count(ctx, channelID)
	if err != nil {
		return nil, err
	}

	textCount, err := s.clipboardRepo.CountByType(ctx, "text", channelID)
	if err != nil {
		return nil, err
	}

	linkCount, err := s.clipboardRepo.CountByType(ctx, "link", channelID)
	if err != nil {
		return nil, err
	}

	codeCount, err := s.clipboardRepo.CountByType(ctx, "code", channelID)
	if err != nil {
		return nil, err
	}

	passwordCount, err := s.clipboardRepo.CountByType(ctx, "password", channelID)
	if err != nil {
		return nil, err
	}

	onlineDevices, err := s.deviceRepo.CountOnline(ctx, channelID)
	if err != nil {
		return nil, err
	}

	totalDevices, err := s.deviceRepo.CountTotal(ctx, channelID)
	if err != nil {
		return nil, err
	}

	syncCount, err := s.syncEventRepo.Count(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return &service.StatsOutput{
		Clipboard: service.ClipboardStats{
			Total:    clipboardCount,
			Text:     textCount,
			Link:     linkCount,
			Code:     codeCount,
			Password: passwordCount,
		},
		Devices: service.DevicesStats{
			Online: onlineDevices,
			Total:  totalDevices,
		},
		SyncCount: syncCount,
	}, nil
}
