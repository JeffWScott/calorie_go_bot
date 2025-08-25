package server

import (
	"calorie_bot/types"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const MODEL_NAME_DEFAULT = "x-ai/grok-4"

func (s *Server) createMeal(c *fiber.Ctx) error {
	var body struct {
		Prompt string  `json:"prompt"`
		Model  *string `json:"model"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	if body.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "prompt parameter is required"})
	}

	model := MODEL_NAME_DEFAULT
	if body.Model != nil {
		model = *body.Model
	}

	job := types.Job{
		Prompt: body.Prompt,
		Model:  model,
		Id:     "1234",
	}

	response, err := s.pipeline.Run(job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": response})
}
