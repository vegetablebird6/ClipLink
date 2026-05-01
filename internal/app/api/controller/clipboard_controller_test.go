package controller

import "testing"

func TestKeysetHasMoreTrimsExtraRecord(t *testing.T) {
	items, hasMore := keysetHasMore([]string{"newest", "middle", "extra"}, 2)

	if !hasMore {
		t.Fatalf("expected hasMore=true when repository returns size+1 records")
	}
	if len(items) != 2 {
		t.Fatalf("expected trimmed page size 2, got %d", len(items))
	}
	if items[0] != "newest" || items[1] != "middle" {
		t.Fatalf("unexpected trimmed items: %#v", items)
	}
}

func TestKeysetHasMoreAllowsExactPageWithoutMore(t *testing.T) {
	items, hasMore := keysetHasMore([]string{"newest", "oldest"}, 2)

	if hasMore {
		t.Fatalf("expected hasMore=false for exact page without extra record")
	}
	if len(items) != 2 {
		t.Fatalf("expected page to remain size 2, got %d", len(items))
	}
}
