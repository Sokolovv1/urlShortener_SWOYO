package controller

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"urlShortener/internal/model"
	"urlShortener/internal/repository"
	"urlShortener/internal/service"
)

type ShortenerController struct {
	shortenerService service.ShortenerServiceInterface
	logger           *zap.Logger
}

func NewShortenerController(svc service.ShortenerServiceInterface, logger *zap.Logger) *ShortenerController {
	return &ShortenerController{
		shortenerService: svc,
		logger:           logger,
	}
}

func (s *ShortenerController) CreateShortenerURL(c fiber.Ctx) error {
	var req model.Request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "URL must not be nil"})
	}

	resp, err := s.shortenerService.CreateShortURL(c.Context(), req.URL)
	if err != nil {
		s.logger.Error("Failed to create short url", zap.String("url", req.URL), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *ShortenerController) GetOriginalURL(c fiber.Ctx) error {

	shortenerURL := c.Params("shortenerURL")

	resp, err := s.shortenerService.GetOriginalURL(c.Context(), shortenerURL)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": repository.ErrLinkNotFound.Error()})
		}
		s.logger.Error("some error occurred", zap.String("shortenerURL", shortenerURL), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
