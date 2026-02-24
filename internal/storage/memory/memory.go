package memory

import (
	"context"
	"fmt"
	"linkShorter/internal/shorter"
)

type MemoryStorage struct {
	urls map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

func (ms *MemoryStorage) SaveUrl(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url cannot be empty")
	}
	if _, ok := ms.urls[url]; ok {
		return "", fmt.Errorf("url already exists")
	}
	short := shorter.GenerateShort()
	ms.urls[short] = url
	return short, nil
}

func (ms *MemoryStorage) GetUrl(ctx context.Context, short string) (string, error) {
	url, ok := ms.urls[short]
	if !ok {
		return "", fmt.Errorf("url not found")
	}
	return url, nil
}
