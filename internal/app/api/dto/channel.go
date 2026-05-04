package dto

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// --- Request DTOs ---

// CreateChannelRequest 创建频道请求
type CreateChannelRequest struct {
	ChannelID string `json:"channel_id"`
}

// VerifyChannelRequest 验证频道请求
type VerifyChannelRequest struct {
	ChannelID string `json:"channel_id"`
}

// --- Response DTOs ---

// ChannelResponse 频道响应
type ChannelResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// ChannelDeleteResponse 频道删除结果响应
type ChannelDeleteResponse struct {
	ChannelID             string `json:"channel_id"`
	ClipboardItemsDeleted int64  `json:"clipboard_items_deleted"`
	SyncEventsDeleted     int64  `json:"sync_events_deleted"`
	DeviceLinksDeleted    int64  `json:"device_links_deleted"`
	OrphanDevicesDeleted  int64  `json:"orphan_devices_deleted"`
}

// --- Converters ---

// ToChannelResponse 从 model.Channel 创建响应
func ToChannelResponse(channel *model.Channel) *ChannelResponse {
	return &ChannelResponse{
		ID:        channel.ID,
		CreatedAt: channel.CreatedAt,
	}
}

// ToChannelDeleteResponse 从 model.ChannelDeleteResult 创建响应
func ToChannelDeleteResponse(r *model.ChannelDeleteResult) *ChannelDeleteResponse {
	return &ChannelDeleteResponse{
		ChannelID:             r.ChannelID,
		ClipboardItemsDeleted: r.ClipboardItemsDeleted,
		SyncEventsDeleted:     r.SyncEventsDeleted,
		DeviceLinksDeleted:    r.DeviceLinksDeleted,
		OrphanDevicesDeleted:  r.OrphanDevicesDeleted,
	}
}
