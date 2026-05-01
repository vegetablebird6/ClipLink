package persistence

import (
	"testing"
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
