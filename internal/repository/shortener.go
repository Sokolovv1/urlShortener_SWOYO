package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"urlShortener/internal/initialize"
)

type SwapRepository interface {
	CreateShortURL(ctx context.Context, id int, shortURL string, originalURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	CheckDublicate(ctx context.Context, originalURL string) (string, error)
	GetNextID(ctx context.Context) (int, error)
}

type PgxIface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ShortenerRepository struct {
	pool   PgxIface
	err    error
	logger *zap.Logger
}

func NewShortenerRepository(dbInstance *initialize.DB, logger *zap.Logger) (*ShortenerRepository, error) {

	err := dbInstance.RunMigrations(logger)

	if err != nil {

		return nil, err
	}

	return &ShortenerRepository{
		pool:   dbInstance.Pool,
		logger: logger,
	}, nil
}

func (r *ShortenerRepository) CreateShortURL(ctx context.Context, id int, shortURL string, originalURL string) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO links (id, short_url, original_url) VALUES ($1, $2, $3)", id, shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *ShortenerRepository) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := r.pool.QueryRow(ctx, "SELECT original_url FROM links WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return originalURL, nil
}

func (r *ShortenerRepository) CheckDublicate(ctx context.Context, originalURL string) (string, error) {
	var dublicateURL string
	err := r.pool.QueryRow(ctx, "SELECT short_url FROM links WHERE original_url = $1", originalURL).Scan(&dublicateURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return dublicateURL, nil
}

func (r *ShortenerRepository) GetNextID(ctx context.Context) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, "SELECT COALESCE(MAX(id), 0) + 1 FROM links").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
