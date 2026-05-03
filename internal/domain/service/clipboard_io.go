package service

import "time"

// --- Clipboard Usecase Input ---

// CreateClipboardInput 创建剪贴板条目的输入
type CreateClipboardInput struct {
	ChannelID       string
	ActorDeviceID   string
	ActorDeviceType string
	Title           string
	Content         string
	Type            string
	CleanDuplicates bool
	ContentHTML     string
	ContentFormat   string
}

// UpdateClipboardInput 更新剪贴板条目的输入（nil 字段不更新）
type UpdateClipboardInput struct {
	ID             string
	ChannelID      string
	ActorDeviceID  string
	Title          *string
	Content        *string
	Type           *string
	DeviceType     *string
	ContentHTML    *string
	ContentFormat  *string
}

// SetFavoriteInput 设置收藏状态的输入
type SetFavoriteInput struct {
	ID            string
	ChannelID     string
	ActorDeviceID string
	Favorite      bool
}

// DeleteClipboardInput 删除剪贴板条目的输入
type DeleteClipboardInput struct {
	ID            string
	ChannelID     string
	ActorDeviceID string
}

// --- Clipboard Usecase Output ---

// ClipboardItemOutput 剪贴板条目输出（只暴露 API 需要的字段）
type ClipboardItemOutput struct {
	ID            string
	ChannelID     string
	Content       string
	ContentHTML   string
	ContentFormat string
	Type          string
	Title         string
	DeviceID      string
	DeviceType    string
	Favorite      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
