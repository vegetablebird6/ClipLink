package service

import "time"

// DeviceOutput 设备输出（只暴露 API 需要的字段）
type DeviceOutput struct {
	ID        string
	Name      string
	Type      string
	LastSeen  time.Time
	IsOnline  bool
	CreatedAt time.Time
}

// DeviceChannelOutput 通道内设备输出
type DeviceChannelOutput struct {
	ID        string
	Name      string
	Type      string
	ChannelID string
	LastSeen  time.Time
	IsOnline  bool
	CreatedAt time.Time
	JoinedAt  time.Time
}
