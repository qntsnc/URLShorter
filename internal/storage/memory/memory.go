package memory

import (
	"context"
	"fmt"
	"linkShorter/internal/service/shorter"
	"sync"
)

type MemoryStorage struct {
	mu   sync.RWMutex
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
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, ok := ms.urls[url]; ok {
		return "", fmt.Errorf("url already exists")
	}
	short := shorter.GenerateShort(url)
	ms.urls[short] = url
	return short, nil
}

func (ms *MemoryStorage) GetUrl(ctx context.Context, short string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	url, ok := ms.urls[short]
	if !ok {
		return "", fmt.Errorf("url is not exist")
	}
	return url, nil
}
