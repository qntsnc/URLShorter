package storage

import (
	"context"
	"fmt"
	"linkShorter/internal/storage/memory"
	"linkShorter/internal/storage/postgres"
)

type Storage interface {
	SaveUrl(ctx context.Context, url string) (string, error)
	GetUrl(ctx context.Context, short string) (string, error)
}

func NewStorage(t string) (Storage, error) {
	switch t {
	case "memory":
		return memory.NewMemoryStorage(), nil
	case "postgres":
		var postgresStorage Storage
		postgresStorage, err := postgres.NewPostgresStorage()
		if err != nil {
			return nil, err
		}
		return postgresStorage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", t)
	}
}
