package server

import (
	"calorie_bot/api"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	app *fiber.App
	api *api.OpenRouter
}

func New(api *api.OpenRouter) *Server {
	app := fiber.New()

	s := &Server{
		app: app,
		api: api,
	}

	// Middleware
	s.app.Use(logger.New())

	v1 := app.Group("v1")

	v1.Post("/meal", s.createMeal)

	return s
}

func (s *Server) Run(host string, port string) error {
	connectionString := fmt.Sprintf("%s:%s", host, port)

	return s.app.Listen(connectionString)
}
