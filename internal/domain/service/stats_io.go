package service

// StatsOutput 通道统计输出（领域服务 → 控制器）
type StatsOutput struct {
	Clipboard ClipboardStats
	Devices   DevicesStats
	SyncCount int64
}

// ClipboardStats 剪贴板统计
type ClipboardStats struct {
	Total    int64
	Text     int64
	Link     int64
	Code     int64
	Password int64
}

// DevicesStats 设备统计
type DevicesStats struct {
	Online int64
	Total  int64
}
