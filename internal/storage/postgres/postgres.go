package postgres

import (
	"context"
	"linkShorter/internal/service/shorter"
	pgdb "linkShorter/internal/storage/postgres/db"
	"linkShorter/internal/storage/redis"
	"log"
	"time"

	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresStorage struct {
	Queries *pgdb.Queries
	redis   *redis.RedisStorage
}

func (ps *PostgresStorage) SaveUrl(ctx context.Context, url string) (string, error) {
	shrt := shorter.GenerateShort(url)
	set, err := ps.redis.Client.SetNX(ctx, shrt, url, 10*time.Minute).Result()
	if err != nil {
		return "", err
	}
	if !set {
		log.Print("Url is already exist\n")
		return shrt, nil
	}
	sh, err := ps.Queries.SaveURL(ctx, pgdb.SaveURLParams{
		Shorturl: shrt,
		Url:      url,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("URL '%s' already exists", url)
			return shrt, nil
		}
		return "", err

	}
	return sh, nil
}

func (ps *PostgresStorage) GetUrl(ctx context.Context, short string) (string, error) {
	FullURL, err := ps.redis.Client.Get(ctx, short).Result()
	if err == nil {
		return FullURL, nil
	}
	FullURL, err = ps.Queries.GetURL(ctx, short)
	if err != nil {
		return "", err
	}
	return FullURL, nil
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "postgres://server:password@db:5432/linkshorter?sslmode=disable"
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
	}
	queries := pgdb.New(conn)
	redisStorage, err := redis.NewRedisStorage()
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{Queries: queries, redis: redisStorage}, nil

}
