package controller_test

import (
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"io"
	"net/http/httptest"
	"testing"
	"urlShortener/internal/controller"
	"urlShortener/internal/model"
	"urlShortener/internal/repository"
	mockService "urlShortener/mocks"
)

func TestCreateShortURL(t *testing.T) {
	logger, _ := zap.NewProduction()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShortenerService := mockService.NewMockShortenerServiceInterface(ctrl)

	app := fiber.New()

	shortenerController := controller.NewShortenerController(mockShortenerService, logger)

	app.Post("/", shortenerController.CreateShortenerURL)

	t.Run("Success", func(t *testing.T) {
		req := &model.Request{URL: "http://example.com"}

		mockShortenerService.EXPECT().
			CreateShortURL(gomock.Any(), req.URL).
			Return(&model.Response{URL: "http://short.url/abc123"}, nil)

		// Создаем новый запрос
		reqBody := `{"URL":"http://example.com"}`
		reqst := httptest.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		reqst.Header.Set("Content-Type", "application/json")

		// Создаем новый Fiber контекст
		resp, err := app.Test(reqst, -1) // -1 означает без тайм-аута
		require.NoError(t, err)

		// Проверяем статус ответа
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Читаем тело ответа
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println(string(body))
		assert.JSONEq(t, `{"url":"http://short.url/abc123"}`, string(body))
	})

	// Тест: Невалидный запрос (пустая ссылка)
	t.Run("invalid request payload", func(t *testing.T) {
		reqBody := ``
		reqst := httptest.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		reqst.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(reqst, -1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, `{"error":"Invalid request payload"}`, string(body))
	})

	t.Run("service error", func(t *testing.T) {
		req := model.Request{URL: "https://example.com"}

		mockShortenerService.EXPECT().
			CreateShortURL(gomock.Any(), req.URL).
			Return(nil, errors.New("internal error"))

		reqBody := `{"url":"https://example.com"}`
		reqst := httptest.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		reqst.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(reqst, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, `{"error":"internal error"}`, string(body))
	})

	t.Run("URL must not be nil", func(t *testing.T) {
		reqBody := `{"URL": ""}`

		reqst := httptest.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		reqst.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(reqst, -1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, `{"error":"URL must not be nil"}`, string(body))
	})
}

func TestGetOriginalURL(t *testing.T) {
	logger, _ := zap.NewProduction()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для ShortenerServiceInterface
	mockShortenerService := mockService.NewMockShortenerServiceInterface(ctrl)

	// Создаем Fiber приложение
	app := fiber.New()

	// Инициализируем контроллер с мок сервисом
	shortenerController := controller.NewShortenerController(mockShortenerService, logger)
	app.Get("/:shortenerURL", shortenerController.GetOriginalURL)

	// Тест: Успешное получение оригинальной ссылки
	t.Run("Success", func(t *testing.T) {
		shortenerURL := "abc123"

		// Определяем поведение мока
		mockShortenerService.EXPECT().
			GetOriginalURL(gomock.Any(), shortenerURL).
			Return(&model.Response{URL: "https://example.com"}, nil)

		reqst := httptest.NewRequest("GET", "/abc123", nil)

		resp, err := app.Test(reqst, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, `{"url":"https://example.com"}`, string(body))
	})

	// Тест: Ошибка при несуществующей короткой ссылке
	t.Run("link not found", func(t *testing.T) {
		shortenerURL := "notfound"

		mockShortenerService.EXPECT().
			GetOriginalURL(gomock.Any(), shortenerURL).
			Return(nil, repository.ErrLinkNotFound)

		reqst := httptest.NewRequest("GET", "/notfound", nil)

		resp, err := app.Test(reqst, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, `{"error":"link not found"}`, string(body))
	})

	// Тест: Ошибка сервиса при получении ссылки
	t.Run("service error", func(t *testing.T) {
		shortenerURL := "abc123"

		mockShortenerService.EXPECT().
			GetOriginalURL(gomock.Any(), shortenerURL).
			Return(nil, errors.New("internal error"))

		reqst := httptest.NewRequest("GET", "/abc123", nil)

		resp, err := app.Test(reqst, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, `{"error":"internal error"}`, string(body))
	})
}
