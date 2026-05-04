package repository

import (
	"context"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// ChannelRepository 通道仓库接口
type ChannelRepository interface {
	// Save 保存通道
	Save(ctx context.Context, channel *model.Channel) error

	// FindByID 通过ID查找通道
	FindByID(ctx context.Context, channelID string) (*model.Channel, error)

	// Exists 检查通道是否存在
	Exists(ctx context.Context, channelID string) (bool, error)

	// Delete 删除通道及其通道内数据，并清理超过指定时间的孤儿设备
	Delete(ctx context.Context, channelID string, orphanDeviceOlderThan time.Time) (*model.ChannelDeleteResult, error)
}
