package dto

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// ChannelStatsResponse 通道统计 API 响应
type ChannelStatsResponse struct {
	TotalDevices       int64     `json:"total_devices"`
	OnlineDevices      int64     `json:"online_devices"`
	ClipboardItemCount int64     `json:"clipboard_item_count"`
	SyncCount          int64     `json:"sync_count"`
	CreatedAt          time.Time `json:"created_at"`
}

// NewChannelStatsResponse 从 StatsOutput 和通道创建时间组装响应
func NewChannelStatsResponse(o *service.StatsOutput, channelCreatedAt time.Time) *ChannelStatsResponse {
	return &ChannelStatsResponse{
		TotalDevices:       o.Devices.Total,
		OnlineDevices:      o.Devices.Online,
		ClipboardItemCount: o.Clipboard.Total,
		SyncCount:          o.SyncCount,
		CreatedAt:          channelCreatedAt,
	}
}
