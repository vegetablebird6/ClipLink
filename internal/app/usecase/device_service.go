package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/xiaojiu/cliplink/internal/common/validation"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"github.com/xiaojiu/cliplink/internal/domain/service"
)

type deviceService struct {
	deviceRepo repository.DeviceRepository
}

func NewDeviceService(deviceRepo repository.DeviceRepository) service.DeviceService {
	return &deviceService{
		deviceRepo: deviceRepo,
	}
}

func (s *deviceService) RegisterDevice(ctx context.Context, name, deviceType, deviceID string) (*service.DeviceOutput, error) {
	if !validation.IsValidDeviceType(deviceType) {
		return nil, model.ErrInvalidInput
	}

	if deviceID == "" {
		deviceID = uuid.New().String()
	}

	existingDevice, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err == nil && existingDevice != nil {
		now := time.Now()
		updates := newDevicePatch().
			withName(name).
			withType(deviceType).
			withLastSeen(now).
			withIsOnline(true).
			toMap()

		if err := s.deviceRepo.Update(ctx, deviceID, updates); err != nil {
			return nil, err
		}

		device, err := s.deviceRepo.FindByID(ctx, deviceID)
		if err != nil {
			return nil, err
		}
		return toDeviceOutput(device), nil
	}

	now := time.Now()
	device := &model.Device{
		ID:        deviceID,
		Name:      name,
		Type:      deviceType,
		LastSeen:  now,
		IsOnline:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.deviceRepo.Save(ctx, device); err != nil {
		return nil, err
	}

	return toDeviceOutput(device), nil
}

func (s *deviceService) GetDeviceByID(ctx context.Context, deviceID string) (*service.DeviceOutput, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	return toDeviceOutput(device), nil
}

func (s *deviceService) UpdateDevice(ctx context.Context, deviceID string, name string, deviceType string) (*service.DeviceOutput, error) {
	p := newDevicePatch().withName(name)
	if deviceType != "" {
		p.withType(deviceType)
	}
	updates := p.toMap()

	if err := s.deviceRepo.Update(ctx, deviceID, updates); err != nil {
		return nil, err
	}

	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	return toDeviceOutput(device), nil
}

func (s *deviceService) UpdateDeviceStatus(ctx context.Context, deviceID string, isOnline bool) (*service.DeviceOutput, error) {
	updates := newDevicePatch().
		withIsOnline(isOnline).
		withLastSeen(time.Now()).
		toMap()

	if err := s.deviceRepo.Update(ctx, deviceID, updates); err != nil {
		return nil, err
	}

	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	return toDeviceOutput(device), nil
}

func (s *deviceService) RemoveDevice(ctx context.Context, deviceID string) error {
	return s.deviceRepo.Delete(ctx, deviceID)
}

func (s *deviceService) AddDeviceToChannel(ctx context.Context, deviceID, channelID string) error {
	existing, err := s.deviceRepo.FindDeviceChannelByDeviceAndChannel(ctx, deviceID, channelID)
	if err != nil {
		return err
	}

	if existing != nil {
		updates := newDeviceChannelPatch().
			withIsActive(true).
			withLastSeenAt(time.Now()).
			toMap()
		return s.deviceRepo.UpdateDeviceChannel(ctx, deviceID, channelID, updates)
	}

	now := time.Now()
	deviceChannel := &model.DeviceChannel{
		DeviceID:   deviceID,
		ChannelID:  channelID,
		IsActive:   true,
		JoinedAt:   now,
		LastSeenAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return s.deviceRepo.SaveDeviceChannel(ctx, deviceChannel)
}

func (s *deviceService) RemoveDeviceFromChannel(ctx context.Context, deviceID, channelID string) error {
	return s.deviceRepo.DeleteDeviceChannel(ctx, deviceID, channelID)
}

func (s *deviceService) UpdateDeviceInChannel(ctx context.Context, deviceID, channelID string, isActive bool) error {
	updates := newDeviceChannelPatch().
		withIsActive(isActive).
		withLastSeenAt(time.Now()).
		toMap()
	return s.deviceRepo.UpdateDeviceChannel(ctx, deviceID, channelID, updates)
}

func (s *deviceService) IsDeviceInChannel(ctx context.Context, deviceID, channelID string) (bool, error) {
	return s.deviceRepo.IsDeviceInChannel(ctx, deviceID, channelID)
}

func (s *deviceService) GetDevicesByChannel(ctx context.Context, channelID string) ([]*service.DeviceChannelOutput, error) {
	dtos, err := s.deviceRepo.FindByChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}
	return toDeviceChannelOutputs(dtos), nil
}

func (s *deviceService) GetDeviceInChannel(ctx context.Context, deviceID, channelID string) (*service.DeviceChannelOutput, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, model.ErrDeviceNotFound
	}

	deviceChannel, err := s.deviceRepo.FindDeviceChannelByDeviceAndChannel(ctx, deviceID, channelID)
	if err != nil {
		return nil, err
	}

	if deviceChannel == nil {
		return nil, model.ErrDeviceNotFound
	}

	return &service.DeviceChannelOutput{
		ID:        device.ID,
		Name:      device.Name,
		Type:      device.Type,
		ChannelID: channelID,
		LastSeen:  device.LastSeen,
		IsOnline:  device.IsOnline,
		CreatedAt: device.CreatedAt,
		JoinedAt:  deviceChannel.JoinedAt,
	}, nil
}

func (s *deviceService) CountOnlineDevices(ctx context.Context, channelID string) (int64, error) {
	return s.deviceRepo.CountOnline(ctx, channelID)
}

func (s *deviceService) CountTotalDevices(ctx context.Context, channelID string) (int64, error) {
	return s.deviceRepo.CountTotal(ctx, channelID)
}

// --- model → output converters ---

func toDeviceOutput(device *model.Device) *service.DeviceOutput {
	if device == nil {
		return nil
	}
	return &service.DeviceOutput{
		ID:        device.ID,
		Name:      device.Name,
		Type:      device.Type,
		LastSeen:  device.LastSeen,
		IsOnline:  device.IsOnline,
		CreatedAt: device.CreatedAt,
	}
}

func toDeviceChannelOutputs(dtos []*model.DeviceDTO) []*service.DeviceChannelOutput {
	result := make([]*service.DeviceChannelOutput, 0, len(dtos))
	for _, d := range dtos {
		result = append(result, &service.DeviceChannelOutput{
			ID:        d.ID,
			Name:      d.Name,
			Type:      d.Type,
			ChannelID: d.ChannelID,
			LastSeen:  d.LastSeen,
			IsOnline:  d.IsOnline,
			CreatedAt: d.CreatedAt,
			JoinedAt:  d.JoinedAt,
		})
	}
	return result
}
