package usecase

import (
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// statsService 统计服务实现
type statsService struct {
	deviceRepo    repository.DeviceRepository
	clipboardRepo repository.ClipboardRepository
	channelRepo   repository.ChannelRepository
	syncEventRepo repository.SyncEventRepository
}

// NewStatsService 创建新的统计服务
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

// GetChannelStats 获取通道统计数据
func (s *statsService) GetChannelStats(channelID string) (*service.StatsOutput, error) {
	// 检查通道是否存在
	exists, err := s.channelRepo.Exists(channelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	// 获取剪贴板统计
	clipboardCount, err := s.clipboardRepo.Count(channelID)
	if err != nil {
		return nil, err
	}

	textCount, err := s.clipboardRepo.CountByType("text", channelID)
	if err != nil {
		return nil, err
	}

	linkCount, err := s.clipboardRepo.CountByType("link", channelID)
	if err != nil {
		return nil, err
	}

	codeCount, err := s.clipboardRepo.CountByType("code", channelID)
	if err != nil {
		return nil, err
	}

	passwordCount, err := s.clipboardRepo.CountByType("password", channelID)
	if err != nil {
		return nil, err
	}

	// 获取设备统计
	onlineDevices, err := s.deviceRepo.CountOnline(channelID)
	if err != nil {
		return nil, err
	}

	totalDevices, err := s.deviceRepo.CountTotal(channelID)
	if err != nil {
		return nil, err
	}

	// 获取同步次数
	syncCount, err := s.syncEventRepo.Count(channelID)
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
