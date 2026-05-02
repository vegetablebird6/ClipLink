package validation

import "testing"

func TestIsValidClipboardType(t *testing.T) {
	validTypes := []string{"text", "link", "code", "password"}
	for _, clipboardType := range validTypes {
		if !IsValidClipboardType(clipboardType) {
			t.Fatalf("expected %q to be valid", clipboardType)
		}
	}

	invalidTypes := []string{"", "image", "file", "other"}
	for _, clipboardType := range invalidTypes {
		if IsValidClipboardType(clipboardType) {
			t.Fatalf("expected %q to be invalid", clipboardType)
		}
	}
}
