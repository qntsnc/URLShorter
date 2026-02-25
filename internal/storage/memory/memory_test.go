package memory

import (
	"context"
	"testing"
)

func TestMemoryStorage_SaveUrl(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
	}{
		{
			name: "valid url",
			url:  "https://example.com",
		},
		{
			name:      "empty url",
			url:       "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemoryStorage()
			shrt, err := ms.SaveUrl(context.Background(), tt.url)
			if (err != nil) != tt.expectErr {
				t.Fatalf("SaveUrl() error = %v, wantErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && shrt == "" {
				t.Error("SaveUrl() returned empty short url")
			}
		})
	}
}

func TestMemoryStorage_SaveUrl_Duplicate(t *testing.T) {
	ms := NewMemoryStorage()
	ctx := context.Background()

	_, err := ms.SaveUrl(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("first SaveUrl() error = %v", err)
	}
	_, err = ms.SaveUrl(ctx, "https://example.com")
	if err != nil {
		t.Error("second SaveUrl() expected handle this and return short url")
	}
}

func TestMemoryStorage_GetUrl(t *testing.T) {
	ms := NewMemoryStorage()
	ctx := context.Background()
	short, err := ms.SaveUrl(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("SaveUrl() error = %v", err)
	}
	got, err := ms.GetUrl(ctx, short)
	if err != nil {
		t.Fatalf("GetUrl() error = %v", err)
	}
	if got != "https://example.com" {
		t.Errorf("GetUrl() = %q, want %q", got, "https://example.com")
	}
}

func TestMemoryStorage_GetUrl_NotFound(t *testing.T) {
	ms := NewMemoryStorage()
	_, err := ms.GetUrl(context.Background(), "nonexistent")
	if err == nil {
		t.Error("GetUrl() expected error for nonexistent short url")
	}
}

func TestMemoryStorage_MultipleSaves(t *testing.T) {
	ms := NewMemoryStorage()
	ctx := context.Background()
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
	}
	shorts := make([]string, len(urls))
	for i, url := range urls {
		short, err := ms.SaveUrl(ctx, url)
		if err != nil {
			t.Fatalf("SaveUrl(%q) error = %v", url, err)
		}
		shorts[i] = short
	}
	for i, short := range shorts {
		got, err := ms.GetUrl(ctx, short)
		if err != nil {
			t.Fatalf("GetUrl(%q) error = %v", short, err)
		}
		if got != urls[i] {
			t.Errorf("GetUrl(%q) = %q, want %q", short, got, urls[i])
		}
	}
}
