package repository

import (
	"context"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// DeviceRepository 设备仓库接口
type DeviceRepository interface {
	// 设备基本操作
	Save(ctx context.Context, device *model.Device) error
	FindByID(ctx context.Context, deviceID string) (*model.Device, error)
	Update(ctx context.Context, deviceID string, updates map[string]interface{}) error
	Delete(ctx context.Context, deviceID string) error

	// 通过ID和通道ID查找设备（用于同步历史快照）
	FindByIDAndChannel(ctx context.Context, deviceID, channelID string) (*model.Device, error)

	// 通道相关设备操作
	FindByChannel(ctx context.Context, channelID string) ([]*model.DeviceDTO, error)
	CountOnline(ctx context.Context, channelID string) (int64, error)
	CountTotal(ctx context.Context, channelID string) (int64, error)

	// 设备通道关联操作
	SaveDeviceChannel(ctx context.Context, deviceChannel *model.DeviceChannel) error
	FindDeviceChannelByDeviceAndChannel(ctx context.Context, deviceID, channelID string) (*model.DeviceChannel, error)
	UpdateDeviceChannel(ctx context.Context, deviceID, channelID string, updates map[string]interface{}) error
	DeleteDeviceChannel(ctx context.Context, deviceID, channelID string) error
	IsDeviceInChannel(ctx context.Context, deviceID, channelID string) (bool, error)
}
