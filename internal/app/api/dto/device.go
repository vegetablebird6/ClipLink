package dto

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
)

// --- Request DTOs ---

// RegisterDeviceRequest 注册设备请求
type RegisterDeviceRequest struct {
	DeviceID   string `json:"device_id" binding:"required"`
	DeviceName string `json:"device_name" binding:"required"`
	DeviceType string `json:"device_type" binding:"required"`
}

// UpdateDeviceStatusRequest 更新设备状态请求
type UpdateDeviceStatusRequest struct {
	IsOnline *bool `json:"is_online" binding:"required"`
}

// UpdateDeviceNameRequest 更新设备名称请求
type UpdateDeviceNameRequest struct {
	Name string `json:"device_name" binding:"required"`
}

// --- Response DTOs ---

// DeviceResponse 设备响应
type DeviceResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	ChannelID string    `json:"channel_id"`
	LastSeen  time.Time `json:"last_seen"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at"`
	JoinedAt  time.Time `json:"joined_at"`
}

// --- Converters ---

// ToDeviceResponse 从 model.DeviceDTO 创建响应
func ToDeviceResponse(d *model.DeviceDTO) *DeviceResponse {
	if d == nil {
		return nil
	}
	return &DeviceResponse{
		ID:        d.ID,
		Name:      d.Name,
		Type:      d.Type,
		ChannelID: d.ChannelID,
		LastSeen:  d.LastSeen,
		IsOnline:  d.IsOnline,
		CreatedAt: d.CreatedAt,
		JoinedAt:  d.JoinedAt,
	}
}

// ToDeviceResponseList 批量转换
func ToDeviceResponseList(dtos []*model.DeviceDTO) []*DeviceResponse {
	items := make([]*DeviceResponse, 0, len(dtos))
	for _, d := range dtos {
		items = append(items, ToDeviceResponse(d))
	}
	return items
}
