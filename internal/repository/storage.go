package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"sync"
)

type URLStorage struct {
	mu      sync.Mutex
	storage map[int]string // ID -> Original URL
	shorts  map[string]int // Short URL -> ID
	logger  *zap.Logger
}

func NewURLStorage(logger *zap.Logger) *URLStorage {
	return &URLStorage{
		storage: make(map[int]string),
		shorts:  make(map[string]int),
		logger:  logger,
	}
}

func (s *URLStorage) CreateShortURL(ctx context.Context, id int, shortURL string, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.shorts[shortURL]; exists {
		s.logger.Error("short URL already exists", zap.Int("id", id), zap.String("short_url", shortURL))
		return errors.New("short URL already exists")
	}

	s.storage[id] = originalURL
	s.shorts[shortURL] = id
	s.logger.Info("short URL created", zap.Int("id", id), zap.String("original_url", originalURL), zap.String("short_url", shortURL))
	return nil
}

func (s *URLStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, exists := s.shorts[shortURL]
	if !exists {
		s.logger.Error("short URL not found", zap.String("short_url", shortURL))
		return "", ErrLinkNotFound
	}

	originalURL, exists := s.storage[id]
	if !exists {
		s.logger.Error("short URL not found", zap.String("short_url", shortURL))
		return "", ErrLinkNotFound
	}

	s.logger.Info("Successfully retrieved short URL", zap.String("original_url", originalURL), zap.String("short_url", shortURL))
	return originalURL, nil
}

func (s *URLStorage) CheckDublicate(ctx context.Context, originalURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for shortURL, storedID := range s.shorts {
		if storedURL, exists := s.storage[storedID]; exists && storedURL == originalURL {
			s.logger.Info("Dublicate short URL found", zap.String("original_url", originalURL))
			return shortURL, nil
		}
	}
	s.logger.Info("Dublicate short URL not found", zap.String("original_url", originalURL))
	return "", ErrLinkNotFound
}

func (s *URLStorage) GetNextID(ctx context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the highest ID and return the next one
	var maxID int
	for id := range s.storage {
		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1, nil
}
