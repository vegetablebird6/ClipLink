package persistence

import (
	"testing"
	"time"

	"github.com/xiaojiu/cliplink/internal/config"
	"github.com/xiaojiu/cliplink/internal/domain/model"
	"github.com/xiaojiu/cliplink/internal/infra/db"
)

func TestCollectDuplicatesUsesContentHashKeys(t *testing.T) {
	repo := &clipboardRepository{}
	seen := make(map[string]struct{})
	duplicateIDs := make([]string, 0)

	repo.collectDuplicates([]duplicateCandidate{
		{ID: "keep", ContentHash: "same-hash"},
		{ID: "duplicate", ContentHash: "same-hash"},
		{ID: "empty-hash"},
	}, seen, &duplicateIDs)

	if len(duplicateIDs) != 1 || duplicateIDs[0] != "duplicate" {
		t.Fatalf("expected duplicate hash row only, got %#v", duplicateIDs)
	}
	if _, exists := seen["same-hash"]; !exists {
		t.Fatalf("expected seen map to store content hash")
	}
}

func TestClipboardRepositoryFindWithKeysetUsesSizePlusOne(t *testing.T) {
	database := setupPersistenceTestDB(t)
	defer closePersistenceTestDB(t, database)

	repo := NewClipboardRepository()
	now := time.Now().UTC()
	channelID := "keyset-channel"

	createClipboardFixtures(t, []model.ClipboardItem{
		{ID: "item-04", ChannelID: channelID, Content: "four", Type: model.TypeText, CreatedAt: now.Add(4 * time.Minute), UpdatedAt: now.Add(4 * time.Minute)},
		{ID: "item-03", ChannelID: channelID, Content: "three", Type: model.TypeText, CreatedAt: now.Add(3 * time.Minute), UpdatedAt: now.Add(3 * time.Minute)},
		{ID: "item-02", ChannelID: channelID, Content: "two", Type: model.TypeText, CreatedAt: now.Add(2 * time.Minute), UpdatedAt: now.Add(2 * time.Minute)},
		{ID: "item-01", ChannelID: channelID, Content: "one", Type: model.TypeText, CreatedAt: now.Add(1 * time.Minute), UpdatedAt: now.Add(1 * time.Minute)},
	})

	firstPage, err := repo.FindWithKeyset(channelID, nil, nil, 2)
	if err != nil {
		t.Fatalf("find first page: %v", err)
	}
	assertClipboardIDs(t, firstPage, []string{"item-04", "item-03", "item-02"})

	afterCreatedAt := firstPage[1].CreatedAt
	afterID := firstPage[1].ID
	secondPage, err := repo.FindWithKeyset(channelID, &afterCreatedAt, &afterID, 2)
	if err != nil {
		t.Fatalf("find second page: %v", err)
	}
	assertClipboardIDs(t, secondPage, []string{"item-02", "item-01"})
}

func TestClipboardRepositoryFindByTypeUsesTypeFilterWithKeyset(t *testing.T) {
	database := setupPersistenceTestDB(t)
	defer closePersistenceTestDB(t, database)

	repo := NewClipboardRepository()
	now := time.Now().UTC()
	channelID := "type-channel"

	createClipboardFixtures(t, []model.ClipboardItem{
		{ID: "text-03", ChannelID: channelID, Content: "text three", Type: model.TypeText, CreatedAt: now.Add(5 * time.Minute), UpdatedAt: now.Add(5 * time.Minute)},
		{ID: "link-02", ChannelID: channelID, Content: "https://example.com", Type: model.TypeLink, CreatedAt: now.Add(4 * time.Minute), UpdatedAt: now.Add(4 * time.Minute)},
		{ID: "text-02", ChannelID: channelID, Content: "text two", Type: model.TypeText, CreatedAt: now.Add(3 * time.Minute), UpdatedAt: now.Add(3 * time.Minute)},
		{ID: "code-01", ChannelID: channelID, Content: "fmt.Println()", Type: model.TypeCode, CreatedAt: now.Add(2 * time.Minute), UpdatedAt: now.Add(2 * time.Minute)},
		{ID: "text-01", ChannelID: channelID, Content: "text one", Type: model.TypeText, CreatedAt: now.Add(1 * time.Minute), UpdatedAt: now.Add(1 * time.Minute)},
	})

	firstPage, err := repo.FindByType(model.TypeText, channelID, nil, nil, 2)
	if err != nil {
		t.Fatalf("find text first page: %v", err)
	}
	assertClipboardIDs(t, firstPage, []string{"text-03", "text-02", "text-01"})

	afterCreatedAt := firstPage[1].CreatedAt
	afterID := firstPage[1].ID
	secondPage, err := repo.FindByType(model.TypeText, channelID, &afterCreatedAt, &afterID, 2)
	if err != nil {
		t.Fatalf("find text second page: %v", err)
	}
	assertClipboardIDs(t, secondPage, []string{"text-01"})
}

func setupPersistenceTestDB(t *testing.T) *db.DB {
	t.Helper()

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	database, err := db.InitWithConfig(&config.Config{Log: config.LogConfig{SQL: "silent"}})
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	return database
}

func closePersistenceTestDB(t *testing.T, database *db.DB) {
	t.Helper()

	if err := database.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
}

func createClipboardFixtures(t *testing.T, items []model.ClipboardItem) {
	t.Helper()

	for i := range items {
		item := items[i]
		if err := db.GetDB().Create(&item).Error; err != nil {
			t.Fatalf("create clipboard fixture %s: %v", item.ID, err)
		}
	}
}

func assertClipboardIDs(t *testing.T, items []*model.ClipboardItem, expected []string) {
	t.Helper()

	if len(items) != len(expected) {
		t.Fatalf("expected %d items %v, got %d %#v", len(expected), expected, len(items), items)
	}
	for i, item := range items {
		if item.ID != expected[i] {
			t.Fatalf("expected item %d to be %s, got %s", i, expected[i], item.ID)
		}
	}
}
