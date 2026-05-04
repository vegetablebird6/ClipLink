package service

import (
	"context"
	"time"
)

// SyncService 同步服务接口
type SyncService interface {
	GetSyncHistory(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *uint, limit int) ([]*SyncEventOutput, error)
	LogSyncAction(ctx context.Context, actorDeviceID, channelID, content string) error
}
