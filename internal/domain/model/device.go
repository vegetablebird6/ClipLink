package model

import (
	"time"
)

// Device 设备模型，用于记录设备基本信息
type Device struct {
	ID        string    `json:"id" gorm:"primarykey"` // 设备ID，主键
	Name      string    `json:"name"`                 // 设备名称
	Type      string    `json:"type"`                 // 设备类型（phone, tablet, desktop, other）
	LastSeen  time.Time `json:"last_seen"`            // 最后一次活跃时间
	IsOnline  bool      `json:"is_online"`            // 是否在线
	CreatedAt time.Time `json:"created_at"`           // 首次创建时间
	UpdatedAt time.Time `json:"updated_at"`           // 更新时间
}

// DeviceDTO 设备数据传输对象，包含设备在特定通道的状态
type DeviceDTO struct {
	ID        string    `json:"id"`         // 设备ID
	Name      string    `json:"name"`       // 设备名称
	Type      string    `json:"type"`       // 设备类型
	ChannelID string    `json:"channel_id"` // 关联的通道ID
	LastSeen  time.Time `json:"last_seen"`  // 最后一次活跃时间
	IsOnline  bool      `json:"is_online"`  // 是否在线
	CreatedAt time.Time `json:"created_at"` // 首次创建时间
	JoinedAt  time.Time `json:"joined_at"`  // 加入通道时间
}

// SyncEvent 同步事件模型，用于记录内容同步事件
// 复合索引由 internal/infra/db/indexes.go 的 EnsureIndexes() 管理。
type SyncEvent struct {
	ID         uint      `json:"id" gorm:"primarykey"`         // 自增ID
	ChannelID  string    `json:"channel_id"`                   // 关联的通道ID
	Action     string    `json:"action"`                       // 动作类型（sync, connect, disconnect, update, delete）
	DeviceID   string    `json:"device_id"`                    // 执行设备ID
	TargetType string    `json:"target_type"`                  // 目标实体类型（clipboard / device / channel）
	TargetID   string    `json:"target_id"`                    // 目标实体 ID
	Content    string    `json:"content"`                      // 操作内容，对于sync是同步的内容摘要
	Summary    string    `json:"summary" gorm:"type:text"`     // 人可读描述
	CreatedAt  time.Time `json:"created_at"`                   // 操作时间
}

// TableName 指定表名为 sync_events
func (SyncEvent) TableName() string { return "sync_events" }

// 同步动作类型常量
const (
	ActionSync       = "sync"       // 同步内容
	ActionConnect    = "connect"    // 设备连接
	ActionDisconnect = "disconnect" // 设备断开连接
	ActionUpdate     = "update"     // 更新内容
	ActionDelete     = "delete"     // 删除内容
)

// 同步目标类型常量
const (
	TargetTypeClipboard = "clipboard" // 剪贴板内容
	TargetTypeDevice    = "device"    // 设备
	TargetTypeChannel   = "channel"   // 通道
)
