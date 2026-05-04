package service

import "context"

// DeviceService 设备服务接口
type DeviceService interface {
	RegisterDevice(ctx context.Context, name, deviceType, deviceID string) (*DeviceOutput, error)
	GetDeviceByID(ctx context.Context, deviceID string) (*DeviceOutput, error)
	UpdateDevice(ctx context.Context, deviceID string, name string, deviceType string) (*DeviceOutput, error)
	UpdateDeviceStatus(ctx context.Context, deviceID string, isOnline bool) (*DeviceOutput, error)
	RemoveDevice(ctx context.Context, deviceID string) error

	AddDeviceToChannel(ctx context.Context, deviceID, channelID string) error
	RemoveDeviceFromChannel(ctx context.Context, deviceID, channelID string) error
	UpdateDeviceInChannel(ctx context.Context, deviceID, channelID string, isActive bool) error
	IsDeviceInChannel(ctx context.Context, deviceID, channelID string) (bool, error)

	GetDevicesByChannel(ctx context.Context, channelID string) ([]*DeviceChannelOutput, error)
	GetDeviceInChannel(ctx context.Context, deviceID, channelID string) (*DeviceChannelOutput, error)
	CountOnlineDevices(ctx context.Context, channelID string) (int64, error)
	CountTotalDevices(ctx context.Context, channelID string) (int64, error)
}
