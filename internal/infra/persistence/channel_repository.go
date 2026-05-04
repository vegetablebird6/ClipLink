package persistence

import (
	"context"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/domain/repository"
	"gorm.io/gorm"
)

// channelRepository 通道仓库实现
type channelRepository struct {
	gdb *gorm.DB
}

// NewChannelRepository 创建新的通道仓库
func NewChannelRepository(gdb *gorm.DB) repository.ChannelRepository {
	return &channelRepository{gdb: gdb}
}

// Save 保存通道
func (r *channelRepository) Save(ctx context.Context, channel *model.Channel) error {
	return r.gdb.WithContext(ctx).Create(channel).Error
}

// FindByID 通过ID查找通道
func (r *channelRepository) FindByID(ctx context.Context, channelID string) (*model.Channel, error) {
	var channel model.Channel
	err := r.gdb.WithContext(ctx).Where("id = ?", channelID).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// Exists 检查通道是否存在
func (r *channelRepository) Exists(ctx context.Context, channelID string) (bool, error) {
	var count int64
	err := r.gdb.WithContext(ctx).Model(&model.Channel{}).Where("id = ?", channelID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Delete 删除通道及其通道内数据，并清理超过指定时间的孤儿设备。
func (r *channelRepository) Delete(ctx context.Context, channelID string, orphanDeviceOlderThan time.Time) (*model.ChannelDeleteResult, error) {
	result := &model.ChannelDeleteResult{ChannelID: channelID}

	err := r.gdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		clipboardResult := tx.Where("channel_id = ?", channelID).Delete(&model.ClipboardItem{})
		if clipboardResult.Error != nil {
			return clipboardResult.Error
		}
		result.ClipboardItemsDeleted = clipboardResult.RowsAffected

		syncEventsResult := tx.Where("channel_id = ?", channelID).Delete(&model.SyncEvent{})
		if syncEventsResult.Error != nil {
			return syncEventsResult.Error
		}
		result.SyncEventsDeleted = syncEventsResult.RowsAffected

		deviceLinksResult := tx.Where("channel_id = ?", channelID).Delete(&model.DeviceChannel{})
		if deviceLinksResult.Error != nil {
			return deviceLinksResult.Error
		}
		result.DeviceLinksDeleted = deviceLinksResult.RowsAffected

		orphanDevicesResult := tx.
			Where("last_seen < ?", orphanDeviceOlderThan).
			Where("NOT EXISTS (?)",
				tx.Model(&model.DeviceChannel{}).
					Select("1").
					Where("device_channels.device_id = devices.id"),
			).
			Delete(&model.Device{})
		if orphanDevicesResult.Error != nil {
			return orphanDevicesResult.Error
		}
		result.OrphanDevicesDeleted = orphanDevicesResult.RowsAffected

		channelResult := tx.Where("id = ?", channelID).Delete(&model.Channel{})
		if channelResult.Error != nil {
			return channelResult.Error
		}
		if channelResult.RowsAffected == 0 {
			return model.ErrChannelNotFound
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
