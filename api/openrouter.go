package api

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

//go:embed meal_macro_schema.json
var schemaJSON embed.FS

var schemaMap map[string]interface{}

func init() {
	data, err := schemaJSON.ReadFile("meal_macro_schema.json")
	if err != nil {
		log.Fatalf("Failed to read schema.json: %v", err)
	}
	if err := json.Unmarshal(data, &schemaMap); err != nil {
		log.Fatalf("Failed to unmarshal schema.json: %v", err)
	}
}

type JSONSchemaMap map[string]interface{}

func (j JSONSchemaMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(j))
}

type OpenRouter struct {
	apiKey  string
	baseURL string // Optional, can be set to "https://openrouter.ai/api/v1" by default
}

type APIClient struct {
	*openai.Client
}

func New() *OpenRouter {
	OPENROUTER_API_KEY, exists := os.LookupEnv("OPENROUTER_API_KEY")
	if !exists {
		log.Fatalf("OPENROUTER_API_KEY not set")
	}

	return &OpenRouter{
		apiKey:  OPENROUTER_API_KEY,
		baseURL: "https://openrouter.ai/api/v1",
	}
}

func (o *OpenRouter) NewClient() *APIClient {
	config := openai.DefaultConfig(o.apiKey)
	config.BaseURL = o.baseURL
	client := openai.NewClientWithConfig(config)

	return &APIClient{Client: client}
}

func (c *APIClient) GetMealMacros(prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: "x-ai/grok-3-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: `
	You are a profesional calorie and macro nutrients counter. 
	A user will provide you a description of a meal and you will provide them back two things:
	1. an array consisting of each part of the mean with it's macros.
	
	Food Grouping example:
	USER - 'I had two baked potatoes with sour cream and three buttered rolls for lunch'

	✅ CORRECT grouping response:
	response [{..., food_name: "Baked Potatoes"}, {..., food_name: "Buttered Rolls"}]

	✖️ BAD grouping response:
	response [{..., food_name: "potato_1"}, {..., food_name: "sour_cream_on_potato_1"}, {..., food_name: "potato_2"}, {..., food_name: "sour_cream_on_potato_2"}, {..., food_name: "roll_1"}, {..., food_name: "butter_on_roll_1"}, {..., food_name: "roll_2"},  {..., food_name: "butter_on_roll_2"}]

	2. a SHORT (6 word MAX) description of the meal they said they had, being as concise as possible. We don't need to specify the meal (lunch) or that it was a meal. Just focue on describing what the macros describe.
	✅ CORRECT meal_description response:
	response {..., meal_description: "Baked Potatoes and Rolls"}

	✖️ BAD meal_description response:
	response {..., meal_description: "Two Baked Potatoes with Sour Cream with Two Buttered Rolls for Lunch."}
	`,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "structured_output",
				Schema: JSONSchemaMap(schemaMap),
			},
		},
	}

	resp, err := c.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil // Or handle as error if needed
	}

	return resp.Choices[0].Message.Content, nil
}
