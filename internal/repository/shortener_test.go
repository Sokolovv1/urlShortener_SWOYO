package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	defer mockPool.Close()

	repo := ShortenerRepository{pool: mockPool}

	// Случай, успешной записи данных
	mockPool.ExpectExec("INSERT INTO links").
		WithArgs(1, "abc123", "https://example.com").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateShortURL(context.Background(), 1, "abc123", "https://example.com")
	assert.NoError(t, err)

	// Случай, когда ошибка при выполнении запроса
	mockPool.ExpectExec("INSERT INTO links").
		WithArgs(1, "abc123", "https://example.com").
		WillReturnError(fmt.Errorf("database error"))

	err = repo.CreateShortURL(context.Background(), 1, "abc123", "https://example.com")
	assert.Error(t, err)
}

func TestGetShortURL(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	defer mockPool.Close()

	repo := ShortenerRepository{pool: mockPool}

	// Случай, когда данные успешно получены
	mockPool.ExpectQuery("SELECT original_url FROM links WHERE short_url").
		WithArgs("abc123").
		WillReturnRows(pgxmock.NewRows([]string{"original_url"}).AddRow("https://example.com"))

	originalURL, err := repo.GetOriginalURL(context.Background(), "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)

	// Случай, когда URL не найден
	mockPool.ExpectQuery("SELECT original_url FROM links WHERE short_url").
		WithArgs("linkNotFound").
		WillReturnError(ErrLinkNotFound)

	_, err = repo.GetOriginalURL(context.Background(), "linkNotFound")
	assert.ErrorIs(t, err, ErrLinkNotFound)

	// Случай, когда ошибка при выполнении запроса
	mockPool.ExpectQuery("SELECT original_url FROM links WHERE short_url").
		WithArgs("abc123").
		WillReturnError(fmt.Errorf("database error"))

	_, err = repo.GetOriginalURL(context.Background(), "abc123")
	assert.Error(t, err)

}

func TestCheckDublicate(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	defer mockPool.Close()

	repo := ShortenerRepository{pool: mockPool}

	// Ситуация, когда дубликат найден
	mockPool.ExpectQuery("SELECT short_url FROM links WHERE original_url").
		WithArgs("https://example.com").WillReturnRows(pgxmock.NewRows([]string{"original_url"}).
		AddRow("abc123"))

	dublicateURL, err := repo.CheckDublicate(context.Background(), "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", dublicateURL)

	// Ситуация, когда дубликат не найден

	mockPool.ExpectQuery("SELECT short_url FROM links WHERE original_url").
		WithArgs("https://example.com").
		WillReturnError(pgx.ErrNoRows)

	_, err = repo.CheckDublicate(context.Background(), "https://example.com")
	assert.ErrorIs(t, err, ErrLinkNotFound)
}

func TestGetNextID(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	defer mockpool.Close()

	repo := ShortenerRepository{pool: mockpool}

	// Случай, когда ожидаем успешное получение нового ID
	mockpool.ExpectQuery("SELECT COALESCE\\(MAX\\(id\\), 0\\) \\+ 1 FROM links").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))

	id, err := repo.GetNextID(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, id)

	// Случай, когда неудалось получить ID
	mockpool.ExpectQuery("SELECT COALESCE\\(MAX\\(id\\), 0\\) \\+ 1 FROM links").
		WillReturnError(fmt.Errorf("database error"))

	id, err = repo.GetNextID(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 0, id)
}
