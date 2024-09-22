package service

import (
	"context"
	"go.uber.org/zap"
	"urlShortener/internal/initialize"
	"urlShortener/internal/model"
	"urlShortener/internal/repository"
	"urlShortener/internal/utils"
)

//go:generate mockgen -source=shortener.go -destination=../../mocks/shortener_mock.go

type SwapRepository interface {
	CreateShortURL(ctx context.Context, id int, shortURL string, originalURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	CheckDublicate(ctx context.Context, originalURL string) (string, error)
	GetNextID(ctx context.Context) (int, error)
}

type ShortenerServiceInterface interface {
	CreateShortURL(ctx context.Context, url string) (*model.Response, error)
	GetOriginalURL(ctx context.Context, url string) (*model.Response, error)
}

type ShortenerService struct {
	repository repository.SwapRepository
	config     *initialize.Config
	logger     *zap.Logger
}

type Deps struct {
	Repository repository.SwapRepository
	Config     *initialize.Config
	Logger     *zap.Logger
}

func NewShortenerService(deps Deps) *ShortenerService {
	return &ShortenerService{
		repository: deps.Repository,
		config:     deps.Config,
		logger:     deps.Logger,
	}
}

func (s *ShortenerService) CreateShortURL(ctx context.Context, url string) (*model.Response, error) {
	existURL, err := s.repository.CheckDublicate(ctx, url)
	if err != nil && err != repository.ErrLinkNotFound {
		return nil, err
	}

	if existURL != "" {
		return &model.Response{
			URL: "http://" + s.config.HTTPHost + ":" + s.config.HTTPPort + "/" + existURL,
		}, nil
	}

	nextID, err := s.repository.GetNextID(ctx)
	if err != nil {
		s.logger.Error("error getting next id", zap.Error(err))
		return nil, err
	}

	shortURL := utils.GenShort(nextID)
	err = s.repository.CreateShortURL(ctx, nextID, shortURL, url)
	if err != nil {
		s.logger.Error("error creating short url", zap.Error(err))
		return nil, err
	}

	return &model.Response{
		URL: "http://" + s.config.HTTPHost + ":" + s.config.HTTPPort + "/" + shortURL,
	}, nil
}

func (s *ShortenerService) GetOriginalURL(ctx context.Context, url string) (*model.Response, error) {
	originalURL, err := s.repository.GetOriginalURL(ctx, url)
	if err != nil {
		s.logger.Error("error getting original url", zap.Error(err))
		return nil, err
	}
	return &model.Response{
		URL: originalURL,
	}, err
}
