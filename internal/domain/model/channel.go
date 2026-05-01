package model

import (
	"time"
)

// Channel 通道模型
type Channel struct {
	ID          string    `json:"id" gorm:"primarykey"` // 通道ID，自动生成的UUID
	Name        string    `json:"name"`                 // 通道名称，可选
	Description string    `json:"description"`          // 通道描述，可选
	CreatedAt   time.Time `json:"created_at"`           // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`           // 更新时间
}

// ChannelDeleteResult describes the rows removed while deleting a channel.
type ChannelDeleteResult struct {
	ChannelID             string `json:"channel_id"`
	ClipboardItemsDeleted int64  `json:"clipboard_items_deleted"`
	SyncEventsDeleted     int64  `json:"sync_events_deleted"`
	DeviceLinksDeleted    int64  `json:"device_links_deleted"`
	OrphanDevicesDeleted  int64  `json:"orphan_devices_deleted"`
}

// DeviceChannel 设备与通道的关联模型 - 解决一个设备可以属于多个通道的问题
// 复合索引由 internal/infra/db/indexes.go 的 EnsureIndexes() 管理。
type DeviceChannel struct {
	ID         uint      `json:"id" gorm:"primarykey"`   // 自增ID
	DeviceID   string    `json:"device_id"`              // 设备ID
	ChannelID  string    `json:"channel_id"`             // 通道ID
	IsActive   bool      `json:"is_active"`              // 设备是否在此通道活跃
	JoinedAt   time.Time `json:"joined_at"`              // 加入通道时间
	LastSeenAt time.Time `json:"last_seen_at"`           // 最后一次在此通道活跃时间
	CreatedAt  time.Time `json:"created_at"`             // 记录创建时间
	UpdatedAt  time.Time `json:"updated_at"`             // 记录更新时间
}
