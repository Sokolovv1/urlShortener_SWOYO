package controller

import "github.com/gofiber/fiber/v3"

func (s *ShortenerController) Register(router fiber.Router) {
	router.Post("/", s.CreateShortenerURL)
	router.Get("/:shortenerURL", s.GetOriginalURL)

}

func (s *ShortenerController) Name() string {
	return ""
}
