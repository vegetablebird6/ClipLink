package dto

import (
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/service"
)

// --- Request DTOs ---

// CreateClipboardRequest 创建剪贴板条目请求
type CreateClipboardRequest struct {
	Title           string `json:"title"`
	Content         string `json:"content" binding:"required"`
	Type            string `json:"type" binding:"required"`
	DeviceID        string `json:"device_id" binding:"required"`
	DeviceType      string `json:"device_type" binding:"required"`
	CleanDuplicates bool   `json:"clean_duplicates"`
	ContentHTML     string `json:"content_html"`
	ContentFormat   string `json:"content_format"`
}

// UpdateClipboardRequest 更新剪贴板条目请求（nil 字段不更新）
type UpdateClipboardRequest struct {
	Title         *string `json:"title"`
	Content       *string `json:"content"`
	Type          *string `json:"type"`
	DeviceID      string  `json:"device_id" binding:"required"`
	DeviceType    *string `json:"device_type"`
	ContentHTML   *string `json:"content_html"`
	ContentFormat *string `json:"content_format"`
}

// SetFavoriteRequest 设置收藏状态请求
type SetFavoriteRequest struct {
	Favorite bool   `json:"favorite"`
	DeviceID string `json:"device_id" binding:"required"`
}

// DeleteClipboardRequest 删除剪贴板条目请求
type DeleteClipboardRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

// --- Response DTOs ---

// ClipboardItemResponse 剪贴板条目响应
type ClipboardItemResponse struct {
	ID            string    `json:"id"`
	ChannelID     string    `json:"channel_id"`
	Content       string    `json:"content"`
	ContentHTML   string    `json:"content_html,omitempty"`
	ContentFormat string    `json:"content_format"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	DeviceID      string    `json:"device_id"`
	DeviceType    string    `json:"device_type"`
	Favorite      bool      `json:"favorite"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// --- Converters ---

// ToClipboardItemResponse 将 usecase output 转换为 API response DTO
func ToClipboardItemResponse(o *service.ClipboardItemOutput) *ClipboardItemResponse {
	if o == nil {
		return nil
	}
	return &ClipboardItemResponse{
		ID:            o.ID,
		ChannelID:     o.ChannelID,
		Content:       o.Content,
		ContentHTML:   o.ContentHTML,
		ContentFormat: o.ContentFormat,
		Type:          o.Type,
		Title:         o.Title,
		DeviceID:      o.DeviceID,
		DeviceType:    o.DeviceType,
		Favorite:      o.Favorite,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}

// ToClipboardItemResponseList 批量转换
func ToClipboardItemResponseList(outputs []*service.ClipboardItemOutput) []*ClipboardItemResponse {
	items := make([]*ClipboardItemResponse, 0, len(outputs))
	for _, o := range outputs {
		items = append(items, ToClipboardItemResponse(o))
	}
	return items
}
