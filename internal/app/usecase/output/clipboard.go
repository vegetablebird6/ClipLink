package output

import "time"

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
