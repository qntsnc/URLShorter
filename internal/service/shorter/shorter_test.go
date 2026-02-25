package shorter

import (
	"testing"
)

func TestGenerateShort(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple url",
			url:  "https://example.com",
		},
		{
			name: "url with path",
			url:  "https://example.com/some/path",
		},
		{
			name: "url with query",
			url:  "https://example.com?q=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shrt := GenerateShort(tt.url)
			if len(shrt) != 10 {
				t.Errorf("GenerateShort(%q) length = %d, want 10", tt.url, len(shrt))
			}
		})
	}
}

func TestGenerateShort_Deterministic(t *testing.T) {
	url := "https://example.com"
	first := GenerateShort(url)
	second := GenerateShort(url)
	if first != second {
		t.Errorf("GenerateShort is not deterministic: %q != %q", first, second)
	}
}

func TestGenerateShort_Unique(t *testing.T) {
	urls := []string{
		"https://example.com",
		"https://example.org",
		"https://google.com",
		"https://github.com",
	}

	results := make(map[string]string)
	for _, url := range urls {
		short := GenerateShort(url)
		if prev, exists := results[short]; exists {
			t.Errorf("Collision: %q and %q both produce %q", prev, url, short)
		}
		results[short] = url
	}
}

func TestGenerateShort_ValidCharacters(t *testing.T) {
	short := GenerateShort("https://example.com")
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	allowedSet := make(map[byte]bool)
	for i := 0; i < len(allowed); i++ {
		allowedSet[allowed[i]] = true
	}
	for i := 0; i < len(short); i++ {
		if !allowedSet[short[i]] {
			t.Errorf("Invalid character %q at position %d", short[i], i)
		}
	}
}
