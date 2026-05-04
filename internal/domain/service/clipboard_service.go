package service

import (
	"context"
	"time"
)

// ClipboardService 剪贴板服务接口
type ClipboardService interface {
	CreateClipboard(ctx context.Context, in CreateClipboardInput) (*ClipboardItemOutput, error)
	GetLatestClipboard(ctx context.Context, channelID string, limit int) ([]*ClipboardItemOutput, error)
	GetClipboardItem(ctx context.Context, id, channelID string) (*ClipboardItemOutput, error)
	GetClipboardHistory(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)
	DeleteClipboard(ctx context.Context, in DeleteClipboardInput) error
	UpdateClipboard(ctx context.Context, in UpdateClipboardInput) (*ClipboardItemOutput, error)
	SetFavorite(ctx context.Context, in SetFavoriteInput) (*ClipboardItemOutput, error)
	GetFavoriteClipboard(ctx context.Context, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)
	GetClipboardByType(ctx context.Context, contentType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)
	GetClipboardByDeviceType(ctx context.Context, deviceType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)
	GetClipboardByTypeAndDeviceType(ctx context.Context, contentType, deviceType string, channelID string, afterCreatedAt *time.Time, afterID *string, size int) ([]*ClipboardItemOutput, error)
	SearchClipboard(ctx context.Context, keyword, channelID string, page, size int) (items []*ClipboardItemOutput, total int64, totalPages int, err error)
	CleanupDuplicateContents(ctx context.Context, channelID string) (int64, error)
}
