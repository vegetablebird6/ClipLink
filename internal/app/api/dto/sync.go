package dto

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// --- Request DTOs ---

// LogSyncActionRequest 记录同步操作请求
type LogSyncActionRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// --- Response DTOs ---

// SyncEventResponse 同步事件响应
type SyncEventResponse struct {
	ID              uint      `json:"id"`
	ChannelID       string    `json:"channel_id"`
	Action          string    `json:"action"`
	TargetType      string    `json:"target_type"`
	TargetID        string    `json:"target_id"`
	Content         string    `json:"content"`
	Summary         string    `json:"summary"`
	ActorDeviceID   string    `json:"actor_device_id"`
	ActorDeviceName string    `json:"actor_device_name"`
	ActorDeviceType string    `json:"actor_device_type"`
	CreatedAt       time.Time `json:"created_at"`
}

// --- Converters ---

// ToSyncEventResponse 从 model.SyncEvent 创建响应
func ToSyncEventResponse(e *model.SyncEvent) *SyncEventResponse {
	if e == nil {
		return nil
	}
	return &SyncEventResponse{
		ID:              e.ID,
		ChannelID:       e.ChannelID,
		Action:          e.Action,
		TargetType:      e.TargetType,
		TargetID:        e.TargetID,
		Content:         e.Content,
		Summary:         e.Summary,
		ActorDeviceID:   e.ActorDeviceID,
		ActorDeviceName: e.ActorDeviceName,
		ActorDeviceType: e.ActorDeviceType,
		CreatedAt:       e.CreatedAt,
	}
}

// ToSyncEventResponseList 批量转换
func ToSyncEventResponseList(events []*model.SyncEvent) []*SyncEventResponse {
	items := make([]*SyncEventResponse, 0, len(events))
	for _, e := range events {
		items = append(items, ToSyncEventResponse(e))
	}
	return items
}
