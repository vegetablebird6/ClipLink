package service

import "github.com/xiaojiu/cliplink/internal/app/usecase/output"

// StatsService 统计服务接口
type StatsService interface {
	// GetChannelStats 获取通道统计数据
	GetChannelStats(channelID string) (*output.StatsOutput, error)
}
