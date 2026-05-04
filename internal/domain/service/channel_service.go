package service

import (
	"context"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// ChannelService 频道服务接口
type ChannelService interface {
	CreateChannel(ctx context.Context, channelID string) (*ChannelOutput, error)
	GetChannel(ctx context.Context, channelID string) (*ChannelOutput, error)
	ChannelExists(ctx context.Context, channelID string) (bool, error)
	VerifyChannel(ctx context.Context, channelID string) (bool, error)
	GetChannelStats(ctx context.Context, channelID string) (*model.ChannelStats, error)
	DeleteChannel(ctx context.Context, channelID string) (*ChannelDeleteOutput, error)
}
