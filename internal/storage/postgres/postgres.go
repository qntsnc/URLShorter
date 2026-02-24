package postgres

import (
	"context"
	"linkShorter/internal/shorter"
	pgdb "linkShorter/internal/storage/postgres/db"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	Queries *pgdb.Queries
}

func (ps *PostgresStorage) SaveUrl(ctx context.Context, url string) (string, error) {
	sh, err := ps.Queries.SaveURL(ctx, pgdb.SaveURLParams{
		Url:      url,
		Shorturl: shorter.GenerateShort()})
	if err != nil {
		return sh, err
	}
	return sh, nil
}

func (ps *PostgresStorage) GetUrl(ctx context.Context, short string) (string, error) {
	FullURL, err := ps.Queries.GetURL(ctx, short)
	if err != nil {
		return "", err
	}
	return FullURL, nil
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "postgres://server:password@postgres:5432/notes_db?sslmode=disable"
	var conn *pgx.Conn
	var err error
	for i := 0; i < 5; i++ {
		conn, err = pgx.Connect(context.Background(), connStr)
		if err == nil {
			log.Println("Successfully connected to database")
			break
		}
		log.Printf("Failed to connect to database (attempt %d/5): %v\n", i+1, err)
		log.Println("Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Unable to connect to database")
		return &PostgresStorage{}, err
	}
	queries := pgdb.New(conn)
	return &PostgresStorage{Queries: queries}, nil

}
