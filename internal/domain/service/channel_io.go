package service

import "time"

// ChannelOutput 频道输出
type ChannelOutput struct {
	ID        string
	CreatedAt time.Time
}

// ChannelDeleteOutput 频道删除结果输出
type ChannelDeleteOutput struct {
	ChannelID             string
	ClipboardItemsDeleted int64
	SyncEventsDeleted     int64
	DeviceLinksDeleted    int64
	OrphanDevicesDeleted  int64
}
