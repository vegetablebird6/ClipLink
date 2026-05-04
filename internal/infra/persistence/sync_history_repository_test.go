package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/infra/db"
)

func TestSyncEventRepositoryFindByChannelUsesSizePlusOne(t *testing.T) {
	database := setupPersistenceTestDB(t)
	defer closePersistenceTestDB(t, database)

	repo := NewSyncEventRepository(db.GetDB())
	now := time.Now().UTC()
	channelID := "sync-keyset-channel"

	createSyncEventFixtures(t, []model.SyncEvent{
		{ChannelID: channelID, Action: model.ActionSync, Content: "four", CreatedAt: now.Add(4 * time.Minute)},
		{ChannelID: channelID, Action: model.ActionSync, Content: "three", CreatedAt: now.Add(3 * time.Minute)},
		{ChannelID: channelID, Action: model.ActionSync, Content: "two", CreatedAt: now.Add(2 * time.Minute)},
		{ChannelID: channelID, Action: model.ActionSync, Content: "one", CreatedAt: now.Add(1 * time.Minute)},
	})

	ctx := context.Background()
	firstPage, err := repo.FindByChannel(ctx, channelID, nil, nil, 2)
	if err != nil {
		t.Fatalf("find first page: %v", err)
	}
	assertSyncEventContents(t, firstPage, []string{"four", "three", "two"})

	afterCreatedAt := firstPage[1].CreatedAt
	afterID := firstPage[1].ID
	secondPage, err := repo.FindByChannel(ctx, channelID, &afterCreatedAt, &afterID, 2)
	if err != nil {
		t.Fatalf("find second page: %v", err)
	}
	assertSyncEventContents(t, secondPage, []string{"two", "one"})
}

func TestSyncEventRepositoryFindByChannelScopesToChannel(t *testing.T) {
	database := setupPersistenceTestDB(t)
	defer closePersistenceTestDB(t, database)

	repo := NewSyncEventRepository(db.GetDB())
	now := time.Now().UTC()

	createSyncEventFixtures(t, []model.SyncEvent{
		{ChannelID: "target-channel", Action: model.ActionSync, Content: "target two", CreatedAt: now.Add(2 * time.Minute)},
		{ChannelID: "other-channel", Action: model.ActionSync, Content: "other", CreatedAt: now.Add(3 * time.Minute)},
		{ChannelID: "target-channel", Action: model.ActionSync, Content: "target one", CreatedAt: now.Add(1 * time.Minute)},
	})

	events, err := repo.FindByChannel(context.Background(), "target-channel", nil, nil, 10)
	if err != nil {
		t.Fatalf("find target channel: %v", err)
	}
	assertSyncEventContents(t, events, []string{"target two", "target one"})
}

func createSyncEventFixtures(t *testing.T, events []model.SyncEvent) {
	t.Helper()

	for i := range events {
		event := events[i]
		if err := db.GetDB().Create(&event).Error; err != nil {
			t.Fatalf("create sync event fixture %s: %v", event.Content, err)
		}
	}
}

func assertSyncEventContents(t *testing.T, events []*model.SyncEvent, expected []string) {
	t.Helper()

	if len(events) != len(expected) {
		t.Fatalf("expected %d events %v, got %d %#v", len(expected), expected, len(events), events)
	}
	for i, event := range events {
		if event.Content != expected[i] {
			t.Fatalf("expected event %d to be %s, got %s", i, expected[i], event.Content)
		}
	}
}
