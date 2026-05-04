package service

import "context"

// ChannelService 频道服务接口
type ChannelService interface {
	CreateChannel(ctx context.Context, channelID string) (*ChannelOutput, error)
	GetChannel(ctx context.Context, channelID string) (*ChannelOutput, error)
	ChannelExists(ctx context.Context, channelID string) (bool, error)
	VerifyChannel(ctx context.Context, channelID string) (bool, error)
	DeleteChannel(ctx context.Context, channelID string) (*ChannelDeleteOutput, error)
}
