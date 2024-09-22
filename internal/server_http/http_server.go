package http

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Controllers []Controller
	Logger      *zap.Logger
}

type Server struct {
	Controller []Controller
	app        *fiber.App
	Logger     *zap.Logger
}

func NewServer(config ServerConfig) *Server {
	app := fiber.New()

	s := &Server{
		Controller: config.Controllers,
		app:        app,
		Logger:     config.Logger,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start(address string) error {
	if err := s.app.Listen(address); err != nil {
		s.Logger.Error("Error starting server", zap.Error(err))
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) registerRoutes() {
	for _, controller := range s.Controller {
		router := s.app.Group(controller.Name())
		controller.Register(router)
	}
}

func (s *Server) ListRoutes() {
	fmt.Println("Registered routes:")
	for _, stack := range s.app.Stack() {
		for _, route := range stack {
			fmt.Printf("%s %s\n", route.Method, route.Path)
		}
	}
}

type Controller interface {
	Register(router fiber.Router)
	Name() string
}
