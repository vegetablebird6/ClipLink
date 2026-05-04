package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/xiaojiu/cliplink/internal/config"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/infra/db"
)

func TestChannelRepositoryDeleteCascadesChannelDataAndOldOrphanDevices(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	database, err := db.InitWithConfig(&config.Config{Log: config.LogConfig{SQL: "silent"}})
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
	})

	now := time.Now()
	deleteChannelID := "delete-channel"
	otherChannelID := "other-channel"

	fixtures := []interface{}{
		&model.Channel{ID: deleteChannelID, CreatedAt: now, UpdatedAt: now},
		&model.Channel{ID: otherChannelID, CreatedAt: now, UpdatedAt: now},
		&model.ClipboardItem{ID: "clip-delete", ChannelID: deleteChannelID, Content: "delete", Type: model.TypeText, CreatedAt: now, UpdatedAt: now},
		&model.ClipboardItem{ID: "clip-other", ChannelID: otherChannelID, Content: "keep", Type: model.TypeText, CreatedAt: now, UpdatedAt: now},
		&model.SyncEvent{Action: model.ActionSync, Content: "delete", ChannelID: deleteChannelID, CreatedAt: now},
		&model.SyncEvent{Action: model.ActionSync, Content: "keep", ChannelID: otherChannelID, CreatedAt: now},
		&model.Device{ID: "old-linked-delete", LastSeen: now.Add(-45 * 24 * time.Hour), CreatedAt: now, UpdatedAt: now},
		&model.Device{ID: "recent-linked-delete", LastSeen: now.Add(-2 * 24 * time.Hour), CreatedAt: now, UpdatedAt: now},
		&model.Device{ID: "old-linked-other", LastSeen: now.Add(-45 * 24 * time.Hour), CreatedAt: now, UpdatedAt: now},
		&model.Device{ID: "old-orphan", LastSeen: now.Add(-45 * 24 * time.Hour), CreatedAt: now, UpdatedAt: now},
		&model.Device{ID: "recent-orphan", LastSeen: now.Add(-2 * 24 * time.Hour), CreatedAt: now, UpdatedAt: now},
		&model.DeviceChannel{DeviceID: "old-linked-delete", ChannelID: deleteChannelID, IsActive: true, JoinedAt: now, LastSeenAt: now, CreatedAt: now, UpdatedAt: now},
		&model.DeviceChannel{DeviceID: "recent-linked-delete", ChannelID: deleteChannelID, IsActive: true, JoinedAt: now, LastSeenAt: now, CreatedAt: now, UpdatedAt: now},
		&model.DeviceChannel{DeviceID: "old-linked-other", ChannelID: otherChannelID, IsActive: true, JoinedAt: now, LastSeenAt: now, CreatedAt: now, UpdatedAt: now},
	}
	for _, fixture := range fixtures {
		if err := db.GetDB().Create(fixture).Error; err != nil {
			t.Fatalf("create fixture %#v: %v", fixture, err)
		}
	}

	result, err := NewChannelRepository(db.GetDB()).Delete(context.Background(), deleteChannelID, now.Add(-30*24*time.Hour))
	if err != nil {
		t.Fatalf("delete channel: %v", err)
	}

	if result.ClipboardItemsDeleted != 1 || result.SyncEventsDeleted != 1 || result.DeviceLinksDeleted != 2 || result.OrphanDevicesDeleted != 2 {
		t.Fatalf("unexpected delete result: %#v", result)
	}

	assertCount(t, &model.Channel{}, "id = ?", deleteChannelID, 0)
	assertCount(t, &model.Channel{}, "id = ?", otherChannelID, 1)
	assertCount(t, &model.ClipboardItem{}, "channel_id = ?", deleteChannelID, 0)
	assertCount(t, &model.ClipboardItem{}, "channel_id = ?", otherChannelID, 1)
	assertCount(t, &model.SyncEvent{}, "channel_id = ?", deleteChannelID, 0)
	assertCount(t, &model.SyncEvent{}, "channel_id = ?", otherChannelID, 1)
	assertCount(t, &model.DeviceChannel{}, "channel_id = ?", deleteChannelID, 0)
	assertCount(t, &model.DeviceChannel{}, "channel_id = ?", otherChannelID, 1)
	assertCount(t, &model.Device{}, "id = ?", "old-linked-delete", 0)
	assertCount(t, &model.Device{}, "id = ?", "old-orphan", 0)
	assertCount(t, &model.Device{}, "id = ?", "recent-linked-delete", 1)
	assertCount(t, &model.Device{}, "id = ?", "recent-orphan", 1)
	assertCount(t, &model.Device{}, "id = ?", "old-linked-other", 1)
}

func assertCount(t *testing.T, modelValue interface{}, query string, arg interface{}, expected int64) {
	t.Helper()

	var count int64
	if err := db.GetDB().Model(modelValue).Where(query, arg).Count(&count).Error; err != nil {
		t.Fatalf("count %T: %v", modelValue, err)
	}
	if count != expected {
		t.Fatalf("expected %T count %d for %s=%v, got %d", modelValue, expected, query, arg, count)
	}
}
