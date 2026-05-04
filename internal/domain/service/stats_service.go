package service

import "context"

// StatsService 统计服务接口
type StatsService interface {
	GetChannelStats(ctx context.Context, channelID string) (*StatsOutput, error)
}
