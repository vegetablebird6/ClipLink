package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"gorm.io/gorm"
)

// deviceRepository 设备仓库实现
type deviceRepository struct {
	gdb *gorm.DB
}

// NewDeviceRepository 创建新的设备仓库
func NewDeviceRepository(gdb *gorm.DB) repository.DeviceRepository {
	return &deviceRepository{gdb: gdb}
}

// Save 保存设备
func (r *deviceRepository) Save(ctx context.Context, device *model.Device) error {
	return r.gdb.WithContext(ctx).Create(device).Error
}

// FindByID 通过ID查找设备
func (r *deviceRepository) FindByID(ctx context.Context, deviceID string) (*model.Device, error) {
	var device model.Device
	err := r.gdb.WithContext(ctx).Where("id = ?", deviceID).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrDeviceNotFound
		}
		return nil, err
	}
	return &device, nil
}

// FindByIDAndChannel 通过ID和通道ID查找设备
func (r *deviceRepository) FindByIDAndChannel(ctx context.Context, deviceID, channelID string) (*model.Device, error) {
	var device model.Device
	err := r.gdb.WithContext(ctx).
		Joins("JOIN device_channels ON devices.id = device_channels.device_id").
		Where("devices.id = ? AND device_channels.channel_id = ?", deviceID, channelID).
		First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrDeviceNotFound
		}
		return nil, err
	}
	return &device, nil
}

// FindByChannel 查找通道下的所有设备
func (r *deviceRepository) FindByChannel(ctx context.Context, channelID string) ([]*model.DeviceDTO, error) {
	var deviceDTOs []*model.DeviceDTO

	err := r.gdb.WithContext(ctx).Table("devices").
		Select("devices.id, devices.name, devices.type, devices.last_seen, devices.is_online, devices.created_at, "+
			"device_channels.channel_id, device_channels.joined_at").
		Joins("JOIN device_channels ON devices.id = device_channels.device_id").
		Where("device_channels.channel_id = ?", channelID).
		Order("device_channels.last_seen_at DESC").
		Scan(&deviceDTOs).Error

	if err != nil {
		return nil, err
	}
	return deviceDTOs, nil
}

// Update 更新设备
func (r *deviceRepository) Update(ctx context.Context, deviceID string, updates map[string]interface{}) error {
	if updates["updated_at"] == nil {
		updates["updated_at"] = time.Now()
	}

	result := r.gdb.WithContext(ctx).Model(&model.Device{}).
		Where("id = ?", deviceID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrDeviceNotFound
	}

	return nil
}

// Delete 删除设备
func (r *deviceRepository) Delete(ctx context.Context, deviceID string) error {
	if err := r.gdb.WithContext(ctx).Where("device_id = ?", deviceID).Delete(&model.DeviceChannel{}).Error; err != nil {
		return err
	}

	result := r.gdb.WithContext(ctx).Where("id = ?", deviceID).Delete(&model.Device{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return model.ErrDeviceNotFound
	}

	return nil
}

// CountOnline 统计通道下在线设备数量
func (r *deviceRepository) CountOnline(ctx context.Context, channelID string) (int64, error) {
	var count int64
	err := r.gdb.WithContext(ctx).Model(&model.DeviceChannel{}).
		Joins("JOIN devices ON devices.id = device_channels.device_id").
		Where("device_channels.channel_id = ? AND devices.is_online = ?", channelID, true).
		Count(&count).Error
	return count, err
}

// CountTotal 统计通道下所有设备数量
func (r *deviceRepository) CountTotal(ctx context.Context, channelID string) (int64, error) {
	var count int64
	err := r.gdb.WithContext(ctx).Model(&model.DeviceChannel{}).
		Where("channel_id = ?", channelID).
		Count(&count).Error
	return count, err
}

// SaveDeviceChannel 保存设备通道关联
func (r *deviceRepository) SaveDeviceChannel(ctx context.Context, deviceChannel *model.DeviceChannel) error {
	return r.gdb.WithContext(ctx).Create(deviceChannel).Error
}

// FindDeviceChannelByDeviceAndChannel 查找设备通道关联
func (r *deviceRepository) FindDeviceChannelByDeviceAndChannel(ctx context.Context, deviceID, channelID string) (*model.DeviceChannel, error) {
	var deviceChannel model.DeviceChannel
	err := r.gdb.WithContext(ctx).Where("device_id = ? AND channel_id = ?", deviceID, channelID).First(&deviceChannel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &deviceChannel, nil
}

// UpdateDeviceChannel 更新设备通道关联
func (r *deviceRepository) UpdateDeviceChannel(ctx context.Context, deviceID, channelID string, updates map[string]interface{}) error {
	if updates["updated_at"] == nil {
		updates["updated_at"] = time.Now()
	}

	result := r.gdb.WithContext(ctx).Model(&model.DeviceChannel{}).
		Where("device_id = ? AND channel_id = ?", deviceID, channelID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("device channel association not found")
	}

	return nil
}

// DeleteDeviceChannel 删除设备通道关联
func (r *deviceRepository) DeleteDeviceChannel(ctx context.Context, deviceID, channelID string) error {
	result := r.gdb.WithContext(ctx).Where("device_id = ? AND channel_id = ?", deviceID, channelID).
		Delete(&model.DeviceChannel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("device channel association not found")
	}

	return nil
}

// IsDeviceInChannel 检查设备是否在通道中
func (r *deviceRepository) IsDeviceInChannel(ctx context.Context, deviceID, channelID string) (bool, error) {
	var count int64
	err := r.gdb.WithContext(ctx).Model(&model.DeviceChannel{}).
		Where("device_id = ? AND channel_id = ?", deviceID, channelID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
