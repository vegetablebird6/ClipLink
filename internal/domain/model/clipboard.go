package model

import (
	"time"
)

// 内容类型常量
const (
	TypeText     = "text"     // 文本
	TypeLink     = "link"     // 链接
	TypeCode     = "code"     // 代码
	TypePassword = "password" // 密码
	TypeImage    = "image"    // 图片
	TypeFile     = "file"     // 文件
)

// ClipboardItem 剪贴板项目模型
// 复合索引由 internal/infra/db/indexes.go 的 EnsureIndexes() 管理，不使用 gorm index tag（列顺序由 struct 声明顺序决定，不可控）。
type ClipboardItem struct {
	ID            string    `json:"id" gorm:"primarykey"`                              // 唯一标识符
	ChannelID     string    `json:"channel_id"`                                        // 通道ID
	Content       string    `json:"content"`                                           // 内容（纯文本）
	ContentHTML   string    `json:"content_html" gorm:"column:content_html;type:text"` // 富文本HTML内容
	ContentFormat string    `json:"content_format" gorm:"column:content_format;default:plain"` // 内容格式：plain 或 html
	ContentHash   string    `json:"-" gorm:"column:content_hash"`                      // 内容哈希，用于去重。语义：基于 trim 后纯文本 content 计算 SHA-256，同文本不同 HTML 样式视为同一内容，空字符串不参与去重
	Type          string    `json:"type"`                                              // 类型（text, link, code, password, image, file）
	Title         string    `json:"title"`                                             // 标题
	DeviceID      string    `json:"device_id"`                                         // 设备ID
	DeviceType    string    `json:"device_type"`                                       // 设备类型快照（phone, tablet, desktop, other），保存时从 devices.type 复制
	Favorite      bool      `json:"favorite"`                                          // 是否收藏
	CreatedAt     time.Time `json:"created_at"`                                        // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                                        // 更新时间
}
