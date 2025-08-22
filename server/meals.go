package server

import (
	"calorie_bot/types"
	"fmt"

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

	job := types.Job{
		Prompt: body.Prompt,
		Id:     "1234",
	}

	response, err := s.pipeline.Run(job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": response})
}
