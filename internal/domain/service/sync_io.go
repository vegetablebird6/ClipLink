package service

import "time"

// SyncEventOutput 同步事件输出
type SyncEventOutput struct {
	ID              uint
	ChannelID       string
	Action          string
	TargetType      string
	TargetID        string
	Content         string
	Summary         string
	ActorDeviceID   string
	ActorDeviceName string
	ActorDeviceType string
	CreatedAt       time.Time
}
