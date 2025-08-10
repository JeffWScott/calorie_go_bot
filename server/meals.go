package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) createMeal(c *fiber.Ctx) error {
	var body struct {
		Prompt string `json:"prompt"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	if body.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "prompt parameter is required"})
	}

	openRouterClient := s.api.NewClient()

	res, err := openRouterClient.GetMealMacros(body.Prompt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "response not generated"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": res})
}
