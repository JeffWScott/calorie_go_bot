package api

import (
	"bytes"
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
		Model:           "x-ai/grok-4",
		Temperature:     0.5,
		ReasoningEffort: "high",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: `
	You are a professional  calorie and macro nutrients counter. 
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
	
	Report all macros in the numerical representation of the units below.
	Include Estimated Serving Size of the FoodGrouping in grams (g).
	Report Total Fat, Saturated Fat, Trans Fat, Total Carbohydrates, Dietary Fiber, Total Sugars, Added Sugars, Sugar Alcohol, and Protein in grams (g).
	Report Cholesterol and Sodium in milligrams (mg).
	For Vitamins and Minerals, report Vitamin D in micrograms (mcg), and Calcium, Iron, and Potassium in milligrams (mg).
	Report Calories as a numeric value (no unit).

	Example Value response:
	✅ CORRECT numerical representation:
		{..., "Cholestero_mg": 10}  // numerical
	✖️ BAD numerical representation:
		{..., "Cholestero_mg": "10mg"}  // Not numerical
	
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
		log.Printf("Chat completion error: %v", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil // Or handle as error if needed
	}

	content := resp.Choices[0].Message.Content // Assume this is the raw content
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(content)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
